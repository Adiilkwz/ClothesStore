document.addEventListener("DOMContentLoaded", async () => {
    const params = new URLSearchParams(window.location.search);
    const id = params.get("id");

    if (!id) return;

    try {
        const res = await fetch(`/products/${id}`); // Ensure this matches your API route
        if (!res.ok) throw new Error("Product not found");
        const p = await res.json();
        
        const container = document.getElementById("product-container");
        const imgUrl = p.image_url || "https://via.placeholder.com/400";

        let sizeOptions = '<option value="">Select Size</option>';
        if (p.size) {
            const sizes = p.size.split(",").map(s => s.trim());
            sizes.forEach(s => {
                sizeOptions += `<option value="${s}">${s}</option>`;
            });
        } else {
             sizeOptions = '<option value="OneSize">One Size</option>';
        }

        // Logic for Out of Stock
        let stockDisplay = `<span>Stock: ${p.stock_quantity}</span>`;
        let buttonState = "";
        let buttonText = "Add to Cart";
        
        if (p.stock_quantity <= 0) {
            stockDisplay = `<span style="color: red; font-weight: bold;">Out of Stock</span>`;
            buttonState = "disabled style='background-color: #ccc; cursor: not-allowed;'";
            buttonText = "Sold Out";
        }

        container.innerHTML = `
            <div class="product-image">
                <img src="${imgUrl}" alt="${p.name}">
            </div>
            <div class="product-info">
                <h1>${p.name}</h1>
                <p class="price">${p.price_kzt} KZT</p>
                
                <div class="meta">
                    <span>Category: ${p.category}</span> | 
                    ${stockDisplay}
                </div>

                <p class="description">
                    High-quality ${p.name}. Perfect for your collection.
                </p>

                <div class="selection-area">
                    <div class="selection-row">
                        <label for="size-select">Choose Size:</label>
                        <select id="size-select">
                            ${sizeOptions}
                        </select>
                    </div>

                    <div class="selection-row">
                        <label for="qty-input">Quantity:</label>
                        <input type="number" id="qty-input" value="1" min="1" max="${p.stock_quantity}" oninput="checkMax(this, ${p.stock_quantity})">
                    </div>
                </div>

                <button class="btn-large" id="add-btn" ${buttonState} onclick="addToCart(${p.id}, '${p.name}', ${p.price_kzt}, ${p.stock_quantity})">${buttonText}</button>
            </div>
        `;
    } catch (err) {
        console.error(err);
        document.getElementById("product-container").innerHTML = "<p>Error loading product.</p>";
    }
});

// Helper to stop typing numbers larger than max
function checkMax(input, max) {
    if (parseInt(input.value) > max) {
        input.value = max;
        alert(`Sorry, we only have ${max} items in stock.`);
    }
    if (parseInt(input.value) < 1) {
        input.value = 1;
    }
}

// Updated Function with Stock Validation
function addToCart(id, name, price, maxStock) {
    const sizeSelect = document.getElementById("size-select");
    const qtyInput = document.getElementById("qty-input");

    const selectedSize = sizeSelect.value;
    const quantity = parseInt(qtyInput.value);

    // 1. Basic Validation
    if (!selectedSize) {
        alert("Please select a size first!");
        return;
    }
    if (quantity < 1) {
        alert("Quantity must be at least 1");
        return;
    }
    if (quantity > maxStock) {
        alert(`Sorry, you cannot order more than ${maxStock} items.`);
        return;
    }

    let cart = JSON.parse(localStorage.getItem("cart_items") || "[]");

    // 2. Check if item already exists in cart
    const existingIndex = cart.findIndex(item => item.productId === id);

    if (existingIndex > -1) {
        // Calculate potential total
        const currentCartQty = cart[existingIndex].quantity;
        const newTotal = currentCartQty + quantity;

        // 3. STOCK CHECK: (Cart + New) vs Stock
        if (newTotal > maxStock) {
            alert(`Stock Limit Reached!\nYou already have ${currentCartQty} in your cart.\nYou can only add ${maxStock - currentCartQty} more.`);
            return;
        }

        cart[existingIndex].quantity = newTotal;
        cart[existingIndex].size = selectedSize; // Update size if needed
    } else {
        cart.push({
            productId: id,
            name: name,
            price: price,
            size: selectedSize,
            quantity: quantity
        });
    }
    
    localStorage.setItem("cart_items", JSON.stringify(cart));
    
    // Optional: Visual Feedback
    const btn = document.getElementById("add-btn");
    const originalText = btn.innerText;
    btn.innerText = "Added!";
    btn.style.backgroundColor = "#27ae60";
    setTimeout(() => {
        btn.innerText = originalText;
        btn.style.backgroundColor = ""; // Reset
    }, 1000);
}