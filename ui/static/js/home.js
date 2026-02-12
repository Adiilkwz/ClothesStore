let currentCategory = 'All';

document.addEventListener("DOMContentLoaded", async () => {
    loadGrid('/products?limit=4', 'latest-grid');
    applyFilters();
});

async function loadGrid(url, elementId) {
    const grid = document.getElementById(elementId);
    grid.innerHTML = "<p>Loading...</p>";

    try {
        const res = await fetch(url);
        if (!res.ok) throw new Error("Failed to fetch");
        const products = await res.json();
        
        if (!products || products.length === 0) {
            grid.innerHTML = "<p>No products found matching filters.</p>";
            return;
        }

        grid.innerHTML = "";
        products.forEach(p => {
            const img = p.image_url || "https://via.placeholder.com/300";
            grid.innerHTML += `
                <div class="card">
                    <img src="${img}" alt="${p.name}" onerror="this.src='https://via.placeholder.com/300'">
                    <div class="card-body">
                        <h3>${p.name}</h3>
                        <p class="price">${p.price_kzt} KZT</p>
                        <small style="color:#888;">${p.category} | ${p.size}</small><br><br>
                        <a href="/product?id=${p.id}" class="btn">View Details</a>
                    </div>
                </div>
            `;
        });
    } catch (err) { console.error(err); }
}

function setCategory(category, btn) {
    currentCategory = category;

    document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
    btn.classList.add('active');

    applyFilters();
}

function applyFilters() {
    const size = document.getElementById("filter-size").value;
    const minPrice = document.getElementById("filter-min-price").value;
    const maxPrice = document.getElementById("filter-max-price").value;

    let url = `/products?`;
    
    if (currentCategory !== 'All') url += `category=${currentCategory}&`;
    if (size) url += `size=${size}&`;
    if (minPrice) url += `min_price=${minPrice}&`;
    if (maxPrice) url += `max_price=${maxPrice}&`;

    loadGrid(url, 'filter-grid');
}

function clearFilters() {
    document.getElementById("filter-size").value = "";
    document.getElementById("filter-min-price").value = "";
    document.getElementById("filter-max-price").value = "";
    applyFilters();
}