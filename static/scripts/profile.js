document.addEventListener("DOMContentLoaded", function () {
    loadUserProfile();
    loadUserOrders();
    // Обработчики для модальных окон
    document.getElementById("editProfile").addEventListener("click", openEditModal);
    document.getElementById("saveProfile").addEventListener("click", updateProfile);
    document.getElementById("closeModal").addEventListener("click", closeEditModal);
    document.getElementById("closeCarModal").addEventListener("click", closeAddCarModal);
    document.getElementById("saveCar").addEventListener("click", addCar);
    document.getElementById("closeDeleteCarModal").addEventListener("click", closeDeleteCarModal); // Добавлено
    // Обработчик для перехода на страницу поиска машин с GET-запросом
    document.getElementById("goToSearch").addEventListener("click", goToSearchPage);
    // Обработчик для удаления профиля через иконку
    document.getElementById("deleteProfileIcon").addEventListener("click", confirmDeleteProfile);
});

// Функция подтверждения удаления профиля
function confirmDeleteProfile() {
    const isConfirmed = confirm("Вы уверены, что хотите удалить свой профиль? Это действие нельзя отменить.");
    if (!isConfirmed) return;

    const token = localStorage.getItem("jwtToken");
    if (!token) {
        alert("Вы не авторизованы!");
        window.location.href = "/login";
        return;
    }

    try {
        fetch("/api/profile", {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        })
            .then(response => {
                if (!response.ok) throw new Error(`Ошибка ${response.status}: ${response.statusText}`);
                return response.json();
            })
            .then(() => {
                alert("Профиль успешно удален.");
                localStorage.removeItem("jwtToken");
                window.location.href = "/login";
            })
            .catch(error => {
                alert("Ошибка при удалении профиля");
                console.error("Ошибка:", error);
            });
    } catch (error) {
        alert("Ошибка при удалении профиля");
        console.error("Ошибка:", error);
    }
}

// Функция для перехода на страницу поиска машин
function goToSearchPage() {
    const token = localStorage.getItem("jwtToken");

    // Проверяем, авторизован ли пользователь
    if (!token) {
        alert("Вы не авторизованы!");
        window.location.href = "/login"; // Если не авторизованы, переходим на страницу логина
        return;
    }

    // Перенаправляем пользователя на страницу поиска машин
    window.location.href = "/api/search/page";
}

// Функция загрузки профиля пользователя
async function loadUserProfile() {
    const token = localStorage.getItem("jwtToken");

    if (!token) {
        alert("Вы не авторизованы!");
        window.location.href = "/login";  // Переход на страницу логина, если токен не найден
        return;
    }

    try {
        const response = await fetch("/api/user/get", {
            headers: { "Authorization": `Bearer ${token}` }
        });

        if (!response.ok) throw new Error(`Ошибка ${response.status}`);

        const user = await response.json();

        document.getElementById("createdAt").textContent = user.created_at;
        document.getElementById("firstName").textContent = user.first_name || "Не указано";
        document.getElementById("lastName").textContent = user.last_name || "Не указано";
        document.getElementById("email").textContent = user.email;
        document.getElementById("updatedAt").textContent = user.updated_at;
        document.getElementById("role").textContent = user.role;
        document.getElementById("phoneNumber").textContent = user.phone_number;

        if (user.role === "admin") {
            const adminActions = document.getElementById("adminActions");
            adminActions.innerHTML = `
                <button class="action-button" id="addCar">Добавить машину</button>
                <button class="action-button" id="deleteCar">Удалить машину</button>
            `;
            document.getElementById("addCar").addEventListener("click", openAddCarModal);
            document.getElementById("deleteCar").addEventListener("click", openDeleteCarModal);
        }

        document.getElementById("editFirstName").value = user.first_name || "";
        document.getElementById("editLastName").value = user.last_name || "";
        document.getElementById("editPhone").value = user.phone_number || "";

    } catch (error) {
        document.getElementById("errorMessage").style.display = "block";
        console.error("Ошибка:", error);
    }
}

// Функции открытия/закрытия модальных окон
function openEditModal() { document.getElementById("editProfileModal").style.display = "block"; }
function closeEditModal() { document.getElementById("editProfileModal").style.display = "none"; }
function openAddCarModal() { document.getElementById("addCarModal").style.display = "block"; }
function closeAddCarModal() { document.getElementById("addCarModal").style.display = "none"; }
function openDeleteCarModal() { document.getElementById("deleteCarModal").style.display = "block"; }
function closeDeleteCarModal() { document.getElementById("deleteCarModal").style.display = "none"; }

// Функция обновления профиля
async function updateProfile() {
    const firstName = document.getElementById("editFirstName").value.trim();
    const lastName = document.getElementById("editLastName").value.trim();
    const phoneNumber = document.getElementById("editPhone").value.trim();
    const token = localStorage.getItem("jwtToken");

    if (!token) {
        alert("Вы не авторизованы!");
        window.location.href = "/login";  // Переход на страницу логина, если токен не найден
        return;
    }

    try {
        const response = await fetch("/api/user/update", {
            method: "POST",
            headers: { "Content-Type": "application/json", "Authorization": `Bearer ${token}` },
            body: JSON.stringify({ first_name: firstName, last_name: lastName, phone_number: phoneNumber })
        });

        if (!response.ok) throw new Error(`Ошибка ${response.status}: ${await response.text()}`);

        alert("Профиль успешно обновлен!");
        closeEditModal();
        loadUserProfile();

    } catch (error) {
        alert("Ошибка обновления профиля");
        console.error("Ошибка:", error);
    }
}

// Функция добавления машины
async function addCar() {
    const carData = new FormData();

    // Добавляем текстовые данные
    carData.append("brand", document.getElementById("carMake").value.trim());
    carData.append("model", document.getElementById("carModel").value.trim());
    carData.append("year", parseInt(document.getElementById("carYear").value.trim(), 10));
    carData.append("color", document.getElementById("carColor").value.trim());
    carData.append("mileage", parseInt(document.getElementById("carMileage").value.trim(), 10));
    carData.append("price_per_day", parseFloat(document.getElementById("carPrice").value.trim()));
    carData.append("status", document.getElementById("carStatus").value);
    carData.append("location_id", parseInt(document.getElementById("carLocation").value.trim(), 10));

    // Добавляем файл изображения
    const imageFile = document.getElementById("carImage").files[0];
    if (!imageFile) {
        alert("Пожалуйста, выберите изображение!");
        return;
    }
    carData.append("image", imageFile);

    const token = localStorage.getItem("jwtToken");
    if (!token) {
        alert("Вы не авторизованы!");
        window.location.href = "/login";
        return;
    }

    try {
        const response = await fetch("/api/cars/add", {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
            },
            body: carData, // FormData автоматически устанавливает правильный Content-Type
        });

        if (!response.ok) throw new Error(`Ошибка ${response.status}: ${await response.text()}`);

        alert("Машина успешно добавлена!");
        closeAddCarModal();

    } catch (error) {
        alert("Ошибка при добавлении машины");
        console.error("Ошибка:", error);
    }
}

// Функция удаления машины
async function deleteCar() {
    const carId = document.getElementById("carIdToDelete").value.trim();
    const token = localStorage.getItem("jwtToken");

    if (!token) {
        alert("Вы не авторизованы!");
        window.location.href = "/login";
        return;
    }

    if (!carId || isNaN(carId)) {
        alert("Пожалуйста, введите корректный ID машины.");
        return;
    }

    const isConfirmed = confirm(`Вы уверены, что хотите удалить машину с ID ${carId}? Это действие нельзя отменить.`);
    if (!isConfirmed) return;

    try {
        const response = await fetch(`/api/cars/${carId}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) throw new Error(`Ошибка ${response.status}: ${await response.text()}`);

        alert("Машина успешно удалена!");
        closeDeleteCarModal();

    } catch (error) {
        alert("Ошибка при удалении машины");
        console.error("Ошибка:", error);
    }
}

async function loadUserOrders() {
    const token = localStorage.getItem("jwtToken");
    if (!token) {
        alert("Вы не авторизованы!");
        window.location.href = "/login";
        return;
    }

    try {
        const response = await fetch("/api/orders/user", {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) throw new Error(`Ошибка ${response.status}: ${await response.text()}`);

        const orders = await response.json();

        // Очищаем предыдущие заказы
        const ordersList = document.getElementById("ordersList");
        ordersList.innerHTML = "";

        // Если заказов нет, показываем сообщение
        if (orders.length === 0) {
            ordersList.innerHTML = "<p>У вас пока нет заказов.</p>";
            return;
        }

        // Отображаем каждый заказ
        orders.forEach(order => {
            const orderDiv = document.createElement("div");
            orderDiv.className = "order-item";

            orderDiv.innerHTML = `
                <h3>Заказ №${order.id}</h3>
                <p><strong>Машина:</strong> ${order.car_brand} ${order.car_model}</p>
                <p><strong>Дата аренды:</strong> ${order.start_date}</p>
                <p><strong>Дата возврата:</strong> ${order.end_date}</p>
                <p><strong>Стоимость:</strong> ${order.total_cost} ₽</p>
                <p><strong>Статус:</strong> ${order.status}</p>
            `;

            ordersList.appendChild(orderDiv);
        });
    } catch (error) {
        console.error("Ошибка при загрузке заказов:", error);
        alert("Не удалось загрузить заказы. Попробуйте позже.");
    }
}

// Привязываем обработчик к кнопке "Удалить" в модальном окне
document.getElementById("confirmDeleteCar").addEventListener("click", deleteCar);

