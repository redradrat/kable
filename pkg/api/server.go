package api

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/gommon/log"

	"github.com/redradrat/kable/pkg/concepts"

	"github.com/labstack/echo/v4"
	"github.com/redradrat/kable/pkg/errors"
	"github.com/redradrat/kable/pkg/repositories"
)

const (
	ConceptsApiPath     = "/concepts"
	RepositoriesApiPath = "/repositories"
)

type Serv struct{}

func StartUp(bind string) {
	serv := Serv{}
	e := echo.New()
	v1 := e.Group("/v1")
	RegisterHandlersV1(v1, &serv)
	e.Static("/", "kable.v1.yaml")
	e.Logger = log.New("kable-server")
	e.Logger.Fatal(e.Start(bind))
}

func RegisterHandlersV1(e *echo.Group, serv *Serv) {
	e.GET(ConceptsApiPath, serv.GetConcepts)
	e.GET(ConceptsApiPath+"/:id", serv.GetConcept)
	e.GET(ConceptsApiPath+"/:id/render", serv.RenderConcept)
	e.GET(RepositoriesApiPath, serv.GetRepositories)
	e.GET(RepositoriesApiPath+"/:id", serv.GetRepository)
	e.PUT(RepositoriesApiPath+"/:id", serv.PutRepository)
	e.DELETE(RepositoriesApiPath+"/:id", serv.DeleteRepository)
	e.GET(RepositoriesApiPath+"/:id"+ConceptsApiPath, serv.GetRepositoryConcepts)
	e.GET(RepositoriesApiPath+"/:id"+ConceptsApiPath+"/:path", serv.GetRepositoryConcept)
	e.GET(RepositoriesApiPath+"/:id"+ConceptsApiPath+"/:path/render", serv.RenderRepositoryConcept)
}

func (serv Serv) GetRepository(ctx echo.Context) error {
	id := getRepoIdFromContext(ctx)
	repo, err := getRepo(id)
	if err != nil {
		return err
	}
	payload := RepositoryPayload{
		URL:    repo.URL,
		GitRef: repo.GitRef,
	}
	return ctx.JSON(http.StatusOK, payload)
}

func getRepoIdFromContext(ctx echo.Context) string {
	return ctx.Param("id")
}

func getRepo(id string) (*repositories.Repository, error) {
	repo, err := repositories.GetRepository(id)
	if err != nil {
		if err == errors.RepositoryUnknownError {
			return nil, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("repository with id '%s' does not exist", id))
		}
		return nil, err
	}
	return &repo, nil
}

func (serv Serv) PutRepository(ctx echo.Context) error {
	payload := new(RepositoryPayload)
	name := ctx.Param("id")
	if err := ctx.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("payload is invalid"))
	}
	repo := repositories.Repository{
		GitRepository: repositories.GitRepository{
			URL:    payload.URL,
			GitRef: payload.GitRef,
		},
		Name: name,
	}
	addMod, err := repositories.AddRepository(repo)
	if err != nil {
		ctx.Logger().Errorf("unable to add repository: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("unable to add repository: %s", err))
	}
	if err := repositories.UpdateRegistry(addMod); err != nil {
		ctx.Logger().Errorf("error updating registry: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error updating registry: %v", err))
	}

	return ctx.JSON(http.StatusOK, NewMessage("Successfully added repository '%s'", name))
}

func (serv Serv) DeleteRepository(ctx echo.Context) error {
	name := ctx.Param("id")
	removeMod, err := repositories.RemoveRepository(name)
	if err != nil {
		ctx.Logger().Errorf("unable to add repository: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("unable to remove repository: %s", err))
	}
	if err := repositories.UpdateRegistry(removeMod); err != nil {
		ctx.Logger().Errorf("error updating registry: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error updating registry: %v", err))
	}
	return ctx.JSON(http.StatusOK, NewMessage("Successfully removed repository '%s'", name))
}

func (serv Serv) GetRepositories(ctx echo.Context) error {
	repos, err := repositories.ListRepositories()
	if err != nil {
		return err
	}
	payload := NewRepositoriesPayload()
	for _, repo := range repos {
		payload.Repositories[repo.Name] = RepositoryPayload{
			URL:    repo.URL,
			GitRef: repo.GitRef,
		}
	}
	return ctx.JSON(200, payload)
}

func (serv Serv) GetConcepts(ctx echo.Context) error {
	ctx.Logger().Infof("'%s' hit by user-agent => %s [%s]", ctx.Path(), ctx.Request().UserAgent(), ctx.RealIP())
	cpts, err := concepts.ListConcepts()
	if err != nil {
		ctx.Logger().Errorf("unable to list concepts: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("unable to list concepts: %v", err))
	}
	payload := NewConceptsPayload()
	for _, cpt := range cpts {
		concept, err := concepts.GetRepoConcept(cpt)
		if err != nil {
			ctx.Logger().Errorf("unable to get concept '%s': %v", cpt.String(), err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("unable to get concept '%s': %v", cpt.String(), err))
		}

		payload.Concepts[cpt.String()] = constructConceptPayloadFrom(concept)
	}
	return ctx.JSON(http.StatusOK, payload)
}

func (serv Serv) GetConcept(ctx echo.Context) error {
	ctx.Logger().Infof("'%s' hit by user-agent => %s [%s]", ctx.Path(), ctx.Request().UserAgent(), ctx.RealIP())
	id := UnmarshalId(getRepoIdFromContext(ctx))
	if !concepts.IsValidConceptIdentifier(id) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("given id '%s' is not a valid concept identifier ", id))
	}
	payload, err := constructConceptPayloadFromCI(concepts.ConceptIdentifier(id))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, payload)
}

func ConceptInputsPayloadFrom(c concepts.Concept) []ConceptInputsPayload {
	if len(c.Inputs.Mandatory) == 0 && len(c.Inputs.Optional) == 0 {
		return nil
	}
	var inputs []ConceptInputsPayload
	for id, input := range c.Inputs.Mandatory {
		inputs = append(inputs, ConceptInputsPayload{
			ID:        id,
			Type:      input.Type.String(),
			Mandatory: true,
		})
	}
	for id, input := range c.Inputs.Optional {
		inputs = append(inputs, ConceptInputsPayload{
			ID:        id,
			Type:      input.Type.String(),
			Mandatory: false,
		})
	}
	sort.Sort(ByID(inputs))
	return inputs
}

func (serv Serv) GetRepositoryConcepts(ctx echo.Context) error {
	ctx.Logger().Infof("'%s' hit by user-agent => %s [%s]", ctx.Path(), ctx.Request().UserAgent(), ctx.RealIP())
	id := getRepoIdFromContext(ctx)
	_, err := getRepo(id)
	if err != nil {
		ctx.Logger().Errorf("could not get repo: %v", err)
		return err
	}
	cpts, err := concepts.ListConcepts()
	if err != nil {
		ctx.Logger().Errorf("unable to list concepts: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("unable to list concepts: %v", err))
	}
	payload := NewConceptsPayload()
	for _, cpt := range cpts {
		if cpt.Repo() != id {
			continue
		}
		concept, err := concepts.GetRepoConcept(cpt)
		if err != nil {
			ctx.Logger().Errorf("unable to get concept '%s': %v", cpt.Concept(), err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("unable to get concept '%s': %v", cpt.Concept(), err))
		}

		payload.Concepts[cpt.String()] = constructConceptPayloadFrom(concept)
	}
	return ctx.JSON(http.StatusOK, payload)
}

func (serv Serv) GetRepositoryConcept(ctx echo.Context) error {
	ctx.Logger().Infof("'%s' hit by user-agent => %s [%s]", ctx.Path(), ctx.Request().UserAgent(), ctx.RealIP())
	ci, err := getRepositoryConceptIdentifierFromContext(ctx)
	if err != nil {
		return err
	}
	payload, err := constructConceptPayloadFromCI(*ci)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, payload)
}

func (serv Serv) RenderConcept(ctx echo.Context) error {
	ctx.Logger().Infof("'%s' hit by user-agent => %s [%s]", ctx.Path(), ctx.Request().UserAgent(), ctx.RealIP())
	inPayload := new(RenderConceptInputPayload)
	if err := ctx.Bind(inPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("inPayload is invalid"))
	}
	ci, err := getConceptIdentifierFromContext(ctx)
	if err != nil {
		return err
	}

	respPayload, err := renderConcept(*ci, inPayload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, respPayload)
}

func (serv Serv) RenderRepositoryConcept(ctx echo.Context) error {
	ctx.Logger().Infof("'%s' hit by user-agent => %s [%s]", ctx.Path(), ctx.Request().UserAgent(), ctx.RealIP())
	inPayload := new(RenderConceptInputPayload)
	if err := ctx.Bind(inPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("inPayload is invalid"))
	}
	ci, err := getRepositoryConceptIdentifierFromContext(ctx)
	if err != nil {
		return err
	}

	respPayload, err := renderConcept(*ci, inPayload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, respPayload)
}

func renderConcept(ci concepts.ConceptIdentifier, inPayload *RenderConceptInputPayload) (RenderConceptResultPayload, error) {
	rdr, err := concepts.RenderConcept(ci.String(), inPayload.Values, concepts.TargetType(inPayload.TargetType), concepts.RenderOpts{
		Local:           false,
		WriteRenderInfo: true,
		Single:          inPayload.SingleManifest,
	})
	if err != nil {
		return RenderConceptResultPayload{}, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("could not render values with given input: %v", err))
	}

	var manifests []string
	manifestCount := len(rdr.Files)
	for _, file := range rdr.Files {
		manifests = append(manifests, file.String())
	}

	origin, err := concepts.GetConceptOriginFromRepository(ci.Repo())
	if err != nil {
		return RenderConceptResultPayload{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("could not compute origin of concept: %v", err))
	}

	respPayload := RenderConceptResultPayload{
		Manifests:     manifests,
		ManifestCount: manifestCount,
		Origin:        origin,
	}
	return respPayload, nil
}

func constructConceptPayloadFromCI(id concepts.ConceptIdentifier) (ConceptPayload, error) {
	cpt, err := concepts.GetRepoConcept(id)
	if err != nil {
		return ConceptPayload{}, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("concept with identifier '%s' does not exist", id))
	}
	payload := constructConceptPayloadFrom(cpt)
	return payload, nil
}

func getConceptIdentifierFromContext(ctx echo.Context) (*concepts.ConceptIdentifier, error) {
	id := UnmarshalId(ctx.Param("id"))
	var cid *concepts.ConceptIdentifier
	if !concepts.IsValidConceptIdentifier(id) {
		tempCid := concepts.ConceptIdentifier(id)
		cid = &tempCid
	} else {
		return nil, fmt.Errorf("given concept identifier '%s' is invalid", id)
	}

	return cid, nil
}

func getRepositoryConceptIdentifierFromContext(ctx echo.Context) (*concepts.ConceptIdentifier, error) {
	id := UnmarshalId(getRepoIdFromContext(ctx))
	_, err := getRepo(id)
	if err != nil {
		ctx.Logger().Errorf("could not get repo: %v", err)
		return nil, err
	}
	path := strings.ReplaceAll(ctx.Param("path"), "_", "/")
	ci := concepts.NewConceptIdentifier(path, id)
	return &ci, nil
}

func constructConceptPayloadFrom(cpt *concepts.Concept) ConceptPayload {
	return ConceptPayload{
		Type: cpt.Type.String(),
		Metadata: ConceptMetadataPayload{
			Tags: cpt.Meta.Tags,
			Maintainer: ConceptMaintainerPayload{
				MaintainerName:  cpt.Meta.Maintainer.Name,
				MaintainerEmail: cpt.Meta.Maintainer.Email,
			}},
		Inputs: ConceptInputsPayloadFrom(*cpt),
	}
}

func UnmarshalId(id string) string {
	return strings.ReplaceAll(id, "_", "/")
}

func MarshalId(id string) string {
	return strings.ReplaceAll(id, "/", "_")
}
