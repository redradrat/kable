document.addEventListener('turbo:load', () => {
    // Get all "navbar-burger" elements
    const $navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);

    // Check if there are any navbar burgers
    if ($navbarBurgers.length > 0) {

        // Add click events on each of them
        $navbarBurgers.forEach( el => {

            // Register activator
            el.addEventListener('click', () => {

                // Get the target from the "data-target" attribute
                const target = el.dataset.target;
                const $target = document.getElementById(target);

                // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
                el.classList.toggle('is-active');
                $target.classList.toggle('is-active');

                // Get all kiddie items
                const $dropdownItems = Array.prototype.slice.call($target.querySelectorAll('a.navbar-item'), 0);

                // Add a click event on each of them, that closes the parent menu.
                $dropdownItems.forEach( el => {
                    el.addEventListener('click', () => {
                        // Toggle the "is-active" class on the "navbar-menu"
                        $target.classList.toggle('is-active');
                    });
                });
            });
        });
    }

});