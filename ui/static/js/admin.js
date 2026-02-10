let editingID = null;

document.addEventListener("DOMContentLoaded", () => {
    fetchProducts();

    const form = document.getElementById("add-product-form");
    const submitBtn = form.querySelector("button[type='submit']");
    const formTitle = form.querySelector("h3");

    form.addEventListener("submit", async (e) => {
        e.preventDefault();
        
        const errorBox = document.getElementById("form-error");
        errorBox.style.display = "none";

        const name = document.getElementById("p-name").value.trim();
        const price = parseFloat(document.getElementById("p-price").value);
        const image = document.getElementById("p-image").value.trim();
        const category = document.getElementById("p-category").value.trim();
        const stock = parseInt(document.getElementById("p-stock").value);
        
        const sizeCheckboxes = document.querySelectorAll('input[name="size"]:checked');
        let selectedSizes = [];
        sizeCheckboxes.forEach((cb) => selectedSizes.push(cb.value));

        if (!name || !image || !category) return showError("Please fill in all text fields.");
        if (selectedSizes.length === 0) return showError("Please select at least one size.");
        if (isNaN(price) || price <= 0) return showError("Invalid Price");
        if (isNaN(stock) || stock < 0) return showError("Invalid Stock");

        const productData = {
            name: name,
            price_kzt: price,
            image_url: image,
            category: category,
            size: selectedSizes.join(", "), 
            stock_quantity: stock
        };

        const token = localStorage.getItem("token");
        let url = '/api/admin/products';
        let method = 'POST';

        if (editingID) {
            url = `/api/admin/products/${editingID}`;
            method = 'PUT';
        }

        try {
            const res = await fetch(url, {
                method: method,
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer " + token
                },
                body: JSON.stringify(productData)
            });

            if (res.ok) {
                alert(editingID ? "Product Updated!" : "Product Added!");
                resetForm();
                fetchProducts();
            } else {
                const err = await res.json();
                showError("Error: " + (err.error || "Failed"));
            }
        } catch (error) {
            console.error(error);
            showError("Network Error");
        }
    });
});

async function startEdit(id) {
    try {
        const res = await fetch(`/products/${id}`);
        if (!res.ok) throw new Error("Could not load product");
        const p = await res.json();

        document.getElementById("p-name").value = p.name;
        document.getElementById("p-price").value = p.price_kzt;
        document.getElementById("p-image").value = p.image_url;
        document.getElementById("p-category").value = p.category;
        document.getElementById("p-stock").value = p.stock_quantity;

        const sizes = p.size.split(",").map(s => s.trim());
        document.querySelectorAll('input[name="size"]').forEach(cb => cb.checked = false);

        if (p.size) {
            const sizes = p.size.split(",").map(s => s.trim());
            
            document.querySelectorAll('input[name="size"]').forEach(cb => {
                if (sizes.includes(cb.value)) {
                    cb.checked = true;
                }
            });
        }

        editingID = p.id;
        document.querySelector("#add-product-form h3").innerText = "Edit Product (ID: " + p.id + ")";
        const btn = document.querySelector("#add-product-form button");
        btn.innerText = "Update Product";
        btn.style.background = "#e67e22";
        
        window.scrollTo({ top: 0, behavior: 'smooth' });

    } catch (err) {
        console.error(err);
        alert("Error loading product for edit.");
    }
}

function resetForm() {
    editingID = null;
    document.getElementById("add-product-form").reset();
    document.querySelectorAll('input[name="size"]').forEach(cb => cb.checked = false);
    document.querySelector("#add-product-form h3").innerText = "Add New Product";
    const btn = document.querySelector("#add-product-form button");
    btn.innerText = "Add Product";
    btn.style.background = "";
}

async function fetchProducts() {
    try {
        const res = await fetch('/products');
        const products = await res.json();
        
        const tbody = document.getElementById("admin-products-list");
        tbody.innerHTML = "";

        products.forEach((p, index) => {
            const displaySize = p.size ? p.size : "-";
            const displayStock = p.stock_quantity !== undefined ? p.stock_quantity : 0;

            tbody.innerHTML += `
                <tr>
                    <td>${index + 1}</td> <td>${p.name}</td>
                    <td>${p.price_kzt} KZT</td>
                    <td>${displayStock}</td>
                    <td>${displaySize}</td>
                    <td>
                        <button class="btn-edit" onclick="startEdit(${p.id})">Edit</button>
                        <button onclick="deleteProduct(${p.id})">Delete</button>
                    </td>
                </tr>
            `;
        });
    } catch (err) {
        console.error("Error loading products:", err);
    }
}

async function deleteProduct(id) {
    if(!confirm("Are you sure?")) return;
    const token = localStorage.getItem("token");
    await fetch(`/api/admin/products/${id}`, {
        method: "DELETE",
        headers: { "Authorization": "Bearer " + token }
    });
    fetchProducts();
}

function showError(msg) {
    const el = document.getElementById("form-error");
    el.innerText = msg;
    el.style.display = "block";
}