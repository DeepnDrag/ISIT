document.addEventListener("DOMContentLoaded", function () {
    loadUserProfile();

    // Обработчики для модальных окон
    document.getElementById("editProfile").addEventListener("click", openEditModal);
    document.getElementById("saveProfile").addEventListener("click", updateProfile);
    document.getElementById("closeModal").addEventListener("click", closeEditModal);

    document.getElementById("closeCarModal").addEventListener("click", closeAddCarModal);
    document.getElementById("saveCar").addEventListener("click", addCar);

    // Обработчик для перехода на страницу поиска машин с GET-запросом
    document.getElementById("goToSearch").addEventListener("click", goToSearchPage);
});

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
            adminActions.innerHTML = `<button class="action-button" id="addCar">Добавить машину</button>`;
            document.getElementById("addCar").addEventListener("click", openAddCarModal);
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
function openEditModal() { document.getElementById("editModal").style.display = "block"; }
function closeEditModal() { document.getElementById("editModal").style.display = "none"; }
function openAddCarModal() { document.getElementById("addCarModal").style.display = "block"; }
function closeAddCarModal() { document.getElementById("addCarModal").style.display = "none"; }

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
