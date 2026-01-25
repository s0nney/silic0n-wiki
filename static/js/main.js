document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('search-input');
    const searchResults = document.getElementById('search-results');

    if (!searchInput || !searchResults) return;

    let debounceTimer;

    searchInput.addEventListener('input', function() {
        const query = this.value.trim();

        clearTimeout(debounceTimer);

        if (query.length === 0) {
            searchResults.innerHTML = '';
            searchResults.classList.remove('visible');
            return;
        }

        debounceTimer = setTimeout(function() {
            fetch('/api/search?q=' + encodeURIComponent(query))
                .then(response => response.json())
                .then(data => {
                    if (data.length === 0) {
                        searchResults.innerHTML = '<div class="no-results">No articles found</div>';
                    } else {
                        searchResults.innerHTML = data.map(article =>
                            '<a href="/wiki/' + article.slug + '" class="search-result-item">' +
                            '<div class="search-result-content">' +
                            '<span class="search-result-title">' + article.title + '</span>' +
                            '<span class="search-result-description">' + article.description + '</span>' +
                            '</div>' +
                            '</a>'
                        ).join('');
                    }
                    searchResults.classList.add('visible');
                })
                .catch(error => {
                    console.error('Search error:', error);
                    searchResults.innerHTML = '<div class="no-results">Search error</div>';
                    searchResults.classList.add('visible');
                });
        }, 200);
    });

    document.addEventListener('click', function(e) {
        if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
            searchResults.classList.remove('visible');
        }
    });

    searchInput.addEventListener('focus', function() {
        if (this.value.trim().length > 0 && searchResults.innerHTML !== '') {
            searchResults.classList.add('visible');
        }
    });
});
