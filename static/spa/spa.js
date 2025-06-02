// Global state
let currentUser = null;
let currentRoute = '/login';
let cars = [];
let brands = [];
let models = [];

// API base URL
const API_BASE_URL = 'http://localhost:8080';

// Initialize app
document.addEventListener('DOMContentLoaded', function() {
    checkAuth();
    loadBrands();

    // Form event listeners
    document.getElementById('loginForm').addEventListener('submit', handleLogin);
    document.getElementById('profileForm').addEventListener('submit', handleUpdateProfile);
    document.getElementById('addCarForm').addEventListener('submit', handleAddCar);

    // Brand change listener
    document.getElementById('searchBrand').addEventListener('change', handleBrandChange);
});

// Authentication functions
function checkAuth() {
    const token = localStorage.getItem('token');
    if (token) {
        fetchUserProfile(token);
    } else {
        showPage('login-page');
    }
}

async function fetchUserProfile(token) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/user/get`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (response.ok) {
            currentUser = await response.json();
            document.getElementById('userEmail').textContent = currentUser.email;
            document.getElementById('header').classList.remove('hidden');
            navigateTo('/');
            loadUserProfile();
        } else {
            localStorage.removeItem('token');
            showPage('login-page');
        }
    } catch (error) {
        console.error('Error fetching user profile:', error);
        localStorage.removeItem('token');
        showPage('login-page');
    }
}

async function handleLogin(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const loginData = Object.fromEntries(formData);
    try {
        const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(loginData)
        });
        const result = await response.json();
        if (response.ok) {
            localStorage.setItem('token', result.token);
            checkAuth();
        } else {
            showAlert('loginAlert', result.error, 'error');
        }
    } catch (error) {
        showAlert('loginAlert', 'Network error', 'error');
    }
}

function logout() {
    localStorage.removeItem('token');
    currentUser = null;
    document.getElementById('header').classList.add('hidden');
    showPage('login-page');
}

// Navigation
function navigateTo(route) {
    currentRoute = route;
    switch(route) {
        case '/login':
            showPage('login-page');
            break;
        case '/':
            showPage('search-page');
            loadCars();
            break;
        case '/profile':
            showPage('profile-page');
            loadUserOrders();
            break;
        case '/add-car':
            showPage('add-car-page');
            break;
        default:
            if (route.startsWith('/car/')) {
                const carId = route.split('/')[2];
                showCarDetail(carId);
            }
            break;
    }
}

function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => {
        page.classList.remove('active');
    });
    document.getElementById(pageId).classList.add('active');
}

// Car functions
async function loadBrands() {
    try {
        const response = await fetch(`${API_BASE_URL}/api/cars/brands`, {
            headers: getAuthHeaders()
        });
        if (response.ok) {
            brands = await response.json();
            populateBrandSelect();
        }
    } catch (error) {
        console.error('Error loading brands:', error);
    }
}

function populateBrandSelect() {
    const brandSelect = document.getElementById('searchBrand');
    brandSelect.innerHTML = '<option value="">All Brands</option>';
    brands.forEach(brand => {
        const option = document.createElement('option');
        option.value = brand;
        option.textContent = brand;
        brandSelect.appendChild(option);
    });
}

async function handleBrandChange() {
    const selectedBrand = document.getElementById('searchBrand').value;
    const modelSelect = document.getElementById('searchModel');
    modelSelect.innerHTML = '<option value="">All Models</option>';
    if (selectedBrand) {
        try {
            const response = await fetch(`${API_BASE_URL}/api/cars/models?brand=${encodeURIComponent(selectedBrand)}`, {
                headers: getAuthHeaders()
            });
            if (response.ok) {
                const result = await response.json();
                models = result.models || [];
                models.forEach(model => {
                    const option = document.createElement('option');
                    option.value = model;
                    option.textContent = model;
                    modelSelect.appendChild(option);
                });
            }
        } catch (error) {
            console.error('Error loading models:', error);
        }
    }
}

async function searchCars() {
    const searchParams = {
        brand: document.getElementById('searchBrand').value,
        model: document.getElementById('searchModel').value,
        year: document.getElementById('searchYear').value ?
            parseInt(document.getElementById('searchYear').value) : null,
        max_price: document.getElementById('searchMaxPrice').value ?
            parseFloat(document.getElementById('searchMaxPrice').value) : null
    };

    // Remove empty values
    Object.keys(searchParams).forEach(key => {
        if (!searchParams[key]) {
            delete searchParams[key];
        }
    });

    try {
        const queryString = new URLSearchParams(searchParams).toString();
        const response = await fetch(`${API_BASE_URL}/api/cars/filter?${queryString}`, {
            headers: getAuthHeaders()
        });
        if (response.ok) {
            cars = await response.json();
            displayCars();
        }
    } catch (error) {
        console.error('Error searching cars:', error);
    }
}

async function loadCars() {
    // Load all cars by default
    searchCars();
}

function displayCars() {
    const carsGrid = document.getElementById('carsGrid');
    carsGrid.innerHTML = '';
    if (cars.length === 0) {
        carsGrid.innerHTML = '<p>No cars found matching your criteria.</p>';
        return;
    }
    cars.forEach(car => {
        const carCard = createCarCard(car);
        carsGrid.appendChild(carCard);
    });
}

function createCarCard(car) {
    const card = document.createElement('div');
    card.className = 'car-card';
    card.onclick = () => navigateTo(`/car/${car.id}`);
    card.innerHTML = `
        <div class="car-image">
            ${car.image_url ? `<img src="${car.image_url}" alt="${car.brand} ${car.model}" 
                style="width: 100%; height: 100%; object-fit: cover;">` : 'No Image'}
        </div>
        <div class="car-info">
            <div class="car-title">${car.brand} ${car.model}</div>
            <div class="car-details">
                ${car.year} • ${car.color} • ${car.mileage.toLocaleString()} km
            </div>
            <div class="car-price">$${car.price_per_day}/day</div>
        </div>
    `;
    return card;
}

async function showCarDetail(carId) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/cars/${carId}`, {
            headers: getAuthHeaders()
        });
        if (response.ok) {
            const car = await response.json();
            displayCarDetail(car);
            showPage('car-detail-page');
            // Load reviews
            loadCarReviews(carId);
        }
    } catch (error) {
        console.error('Error loading car details:', error);
    }
}

function displayCarDetail(car) {
    const content = document.getElementById('carDetailContent');
    content.innerHTML = `
        <div class="car-detail">
            <div>
                <div class="car-detail-image">
                    ${car.image_url ? `<img src="${car.image_url}" alt="${car.brand} ${car.model}" 
                        style="width: 100%; height: 100%; object-fit: cover; border-radius: 10px;">` : 'No Image'}
                </div>
            </div>
            <div>
                <h1>${car.brand} ${car.model}</h1>
                <div class="car-specs">
                    <div class="spec-item">
                        <span>Year:</span>
                        <span>${car.year}</span>
                    </div>
                    <div class="spec-item">
                        <span>Color:</span>
                        <span>${car.color}</span>
                    </div>
                    <div class="spec-item">
                        <span>Mileage:</span>
                        <span>${car.mileage.toLocaleString()} km</span>
                    </div>
                    <div class="spec-item">
                        <span>Status:</span>
                        <span>${car.status}</span>
                    </div>
                </div>
                <div class="booking-section">
                    <h3>Book This Car</h3>
                    <div class="date-inputs">
                        <div class="form-group">
                            <label for="startDate">Start Date:</label>
                            <input type="date" id="startDate" required>
                        </div>
                        <div class="form-group">
                            <label for="endDate">End Date:</label>
                            <input type="date" id="endDate" required>
                        </div>
                    </div>
                    <div style="font-size: 1.2rem; margin-bottom: 1rem;">
                        <strong>Price: $${car.price_per_day}/day</strong>
                    </div>
                    <button class="btn btn-full" onclick="bookCar(${car.id}, ${car.price_per_day})">Book Now</button>
                </div>
            </div>
        </div>
        <div class="reviews-section">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
                <h3>Reviews</h3>
                <button class="btn" onclick="showAddReviewForm(${car.id})">Add Review</button>
            </div>
            <div id="addReviewForm" class="hidden" style="background: #f8f9fa; padding: 1rem; border-radius: 8px; margin-bottom: 1rem;">
                <div class="form-group">
                    <label for="reviewRating">Rating:</label>
                    <select id="reviewRating">
                        <option value="1">1 Star</option>
                        <option value="2">2 Stars</option>
                        <option value="3">3 Stars</option>
                        <option value="4">4 Stars</option>
                        <option value="5">5 Stars</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="reviewComment">Comment:</label>
                    <textarea id="reviewComment" rows="3" style="width: 100%; padding: 0.5rem; border: 1px solid #ddd; border-radius: 5px;"></textarea>
                </div>
                <button class="btn" onclick="submitReview(${car.id})">Submit Review</button>
            </div>
            <div id="reviewsContainer">
                <!-- Reviews will be loaded here -->
            </div>
        </div>
    `;
}

async function bookCar(carId, pricePerDay) {
    const startDate = document.getElementById('startDate').value;
    const endDate = document.getElementById('endDate').value;
    if (!startDate || !endDate) {
        alert('Please select both start and end dates');
        return;
    }

    const start = new Date(startDate);
    const end = new Date(endDate);
    const days = Math.ceil((end - start) / (1000 * 60 * 60 * 24));
    const totalCost = days * pricePerDay;

    const orderData = {
        car_id: carId,
        start_date: startDate,
        end_date: endDate,
        total_cost: totalCost
    };

    try {
        const response = await fetch(`${API_BASE_URL}/api/orders`, {
            method: 'POST',
            headers: {
                ...getAuthHeaders(),
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(orderData)
        });
        if (response.ok) {
            alert('Car booked successfully!');
            navigateTo('/profile');
        } else {
            const error = await response.json();
            alert('Booking failed: ' + error.error);
        }
    } catch (error) {
        console.error('Error booking car:', error);
        alert('Network error');
    }
}

async function loadCarReviews(carId) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/reviews/${carId}`);
        if (response.ok) {
            const reviews = await response.json();
            displayReviews(reviews);
        }
    } catch (error) {
        console.error('Error loading reviews:', error);
    }
}

function displayReviews(reviews) {
    const container = document.getElementById('reviewsContainer');
    if (reviews.length === 0) {
        container.innerHTML = '<p>No reviews yet.</p>';
        return;
    }
    container.innerHTML = reviews.map(review => `
        <div class="review-card">
            <div class="review-header">
                <strong>User ${review.user_id}</strong>
                <div class="stars">${'★'.repeat(review.rating)}${'☆'.repeat(5 - review.rating)}</div>
            </div>
            <p>${review.comment}</p>
            <small>${new Date(review.created_at).toLocaleDateString()}</small>
        </div>
    `).join('');
}

function showAddReviewForm(carId) {
    document.getElementById('addReviewForm').classList.toggle('hidden');
}

async function submitReview(carId) {
    const rating = parseInt(document.getElementById('reviewRating').value);
    const comment = document.getElementById('reviewComment').value;

    if (!comment.trim()) {
        alert('Please enter a comment');
        return;
    }

    const reviewData = {
        car_id: carId,
        rating: rating,
        comment: comment
    };

    try {
        const response = await fetch(`${API_BASE_URL}/api/reviews`, {
            method: 'POST',
            headers: {
                ...getAuthHeaders(),
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(reviewData)
        });
        if (response.ok) {
            document.getElementById('reviewComment').value = '';
            document.getElementById('addReviewForm').classList.add('hidden');
            loadCarReviews(carId);
        } else {
            const error = await response.json();
            alert('Failed to submit review: ' + error.error);
        }
    } catch (error) {
        console.error('Error submitting review:', error);
        alert('Network error');
    }
}

// Profile functions
function loadUserProfile() {
    if (currentUser) {
        document.getElementById('firstName').value = currentUser.first_name || '';
        document.getElementById('lastName').value = currentUser.last_name || '';
        document.getElementById('phoneNumber').value = currentUser.phone_number || '';
        document.getElementById('role').value = currentUser.role || '';
    }
}

async function handleUpdateProfile(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const profileData = Object.fromEntries(formData);
    try {
        const response = await fetch(`${API_BASE_URL}/api/user/update`, {
            method: 'POST',
            headers: {
                ...getAuthHeaders(),
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(profileData)
        });
        if (response.ok) {
            currentUser = await response.json();
            alert('Profile updated successfully!');
        } else {
            const error = await response.json();
            alert('Update failed: ' + error.error);
        }
    } catch (error) {
        console.error('Error updating profile:', error);
        alert('Network error');
    }
}

async function loadUserOrders() {
    try {
        const response = await fetch(`${API_BASE_URL}/api/orders/user`, {
            headers: getAuthHeaders()
        });
        if (response.ok) {
            const orders = await response.json();
            displayOrders(orders);
        }
    } catch (error) {
        console.error('Error loading orders:', error);
    }
}

function displayOrders(orders) {
    const container = document.getElementById('ordersContainer');
    if (orders.length === 0) {
        container.innerHTML = '<p>No orders found.</p>';
        return;
    }
    container.innerHTML = orders.map(order => `
        <div class="order-card">
            <div class="order-header">
                <strong>Order #${order.id}</strong>
                <span class="status ${order.status}">${order.status}</span>
            </div>
            <p><strong>Car ID:</strong> ${order.car_id}</p>
            <p><strong>Period:</strong> ${order.start_date} to ${order.end_date}</p>
            <p><strong>Total Cost:</strong> $${order.total_cost}</p>
            <p><strong>Created:</strong> ${new Date(order.created_at).toLocaleDateString()}</p>
        </div>
    `).join('');
}

async function deleteAccount() {
    if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
        try {
            const response = await fetch(`${API_BASE_URL}/users`, {
                method: 'DELETE',
                headers: getAuthHeaders()
            });
            if (response.ok) {
                alert('Account deleted successfully');
                logout();
            } else {
                alert('Failed to delete account');
            }
        } catch (error) {
            console.error('Error deleting account:', error);
            alert('Network error');
        }
    }
}

// Add car functions
async function handleAddCar(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    try {
        const response = await fetch(`${API_BASE_URL}/api/cars/add`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: formData
        });
        if (response.ok) {
            alert('Car added successfully!');
            e.target.reset();
            navigateTo('/');
        } else {
            const error = await response.json();
            alert('Failed to add car: ' + error.error);
        }
    } catch (error) {
        console.error('Error adding car:', error);
        alert('Network error');
    }
}

// Utility functions
function getAuthHeaders() {
    const token = localStorage.getItem('token');
    return {
        'Authorization': `Bearer ${token}`
    };
}

function showAlert(containerId, message, type) {
    const container = document.getElementById(containerId);
    container.innerHTML = `<div class="alert alert-${type}">${message}</div>`;
}