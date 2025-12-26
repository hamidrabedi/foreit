/**
 * Persist changelist filters state (collapsed/expanded).
 */
'use strict';
{
    // Init filters.
    let filters = JSON.parse(sessionStorage.getItem('forge.admin.filtersState'));

    if (!filters) {
        filters = {};
    }

    Object.entries(filters).forEach(([key, value]) => {
        const detailElement = document.querySelector(`[data-filter-title='${CSS.escape(key)}']`);

        // Check if the filter is present, it could be from other view.
        if (detailElement) {
            if (value) {
                detailElement.setAttribute('open', '');
            } else {
                detailElement.removeAttribute('open');
            }
        }
    });

    // Save filter state when clicks.
    const details = document.querySelectorAll('details[data-filter-title]');
    details.forEach(detail => {
        detail.addEventListener('toggle', event => {
            filters[`${event.target.dataset.filterTitle}`] = detail.open;
            sessionStorage.setItem('forge.admin.filtersState', JSON.stringify(filters));
        });
    });
}

