document.addEventListener('DOMContentLoaded', () => {
    const filterBrand = document.getElementById('filterBrand');
    const filterModel = document.getElementById('filterModel');
    const filterYearFrom = document.getElementById('filterYearFrom');
    const filterYearTo = document.getElementById('filterYearTo');
    const startDateInput = document.getElementById('startDate');
    const endDateInput = document.getElementById('endDate');
    const carList = document.getElementById('carList');

    if (!filterBrand || !filterModel || !filterYearFrom || !filterYearTo || !startDateInput || !endDateInput || !carList) {
        console.error("Не найдены необходимые элементы DOM");
        return;
    }

    // Заполняем выпадающий список годов
    function populateYears(selectElement) {
        for (let year = 1920; year <= new Date().getFullYear(); year++) {
            const option = document.createElement('option');
            option.value = year;
            option.textContent = year;
            selectElement.appendChild(option);
        }
    }

    // Установка минимальных дат для полей выбора дат
    function setupDatePickers(startDateInput, endDateInput) {
        const today = new Date().toISOString().split('T')[0]; // Текущая дата в формате YYYY-MM-DD
        startDateInput.setAttribute('min', today);
        endDateInput.setAttribute('min', today);
    }

    // Загрузка брендов с сервера
    async function loadBrands() {
        try {
            const token = localStorage.getItem('jwtToken'); // Получаем токен из localStorage
            if (!token) {
                alert('Вы не авторизованы!');
                window.location.href = '/login'; // Перенаправляем на страницу логина
                return;
            }

            const response = await fetch('/api/cars/brands', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`, // Добавляем заголовок авторизации
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                if (response.status === 401) {
                    alert('Сессия истекла. Пожалуйста, войдите снова.');
                    window.location.href = '/login'; // Перенаправляем на страницу логина
                }
                throw new Error('Ошибка при загрузке брендов');
            }

            const brands = await response.json();

            // Заполняем выпадающий список брендов
            filterBrand.innerHTML = '<option value="">Бренд</option>';
            brands.forEach(brand => {
                const option = document.createElement('option');
                option.value = brand;
                option.textContent = brand;
                filterBrand.appendChild(option);
            });
            console.log(brands);
        } catch (error) {
            console.error(error);
            alert('Не удалось загрузить бренды');
        }
    }

    async function loadModels(brand) {
        try {
            const token = localStorage.getItem('jwtToken'); // Получаем токен из localStorage
            if (!token) {
                alert('Вы не авторизованы!');
                window.location.href = '/login'; // Перенаправляем на страницу логина
                return;
            }

            const response = await fetch(`/api/cars/models?brand=${brand}`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                if (response.status === 401) {
                    alert('Сессия истекла. Пожалуйста, войдите снова.');
                    window.location.href = '/login';
                }
                throw new Error('Ошибка при загрузке моделей');
            }

            const data = await response.json();
            console.log(data);

            const models = data.models;

            if (!Array.isArray(models)) {
                throw new Error('Данные моделей не являются массивом');
            }

            filterModel.innerHTML = '<option value="">Модель</option>';
            models.forEach(model => {
                const option = document.createElement('option');
                option.value = model;
                option.textContent = model;
                filterModel.appendChild(option);
            });
        } catch (error) {
            console.error(error);
            alert('Не удалось загрузить модели');
        }
    }

    async function fetchCars() {
        try {
            const brand = filterBrand.value;
            const model = filterModel.value;
            const yearFrom = filterYearFrom.value;
            const yearTo = filterYearTo.value;
            const minPrice = document.getElementById('filterMinPrice').value;
            const maxPrice = document.getElementById('filterMaxPrice').value;
            const startDate = startDateInput.value;
            const endDate = endDateInput.value;

            const params = new URLSearchParams();
            if (brand) params.append('brand', brand);
            if (model) params.append('model', model);
            if (yearFrom) params.append('year_from', yearFrom);
            if (yearTo) params.append('year_to', yearTo);
            if (minPrice) params.append('min_price', minPrice);
            if (maxPrice) params.append('max_price', maxPrice);
            if (startDate) params.append('start_date', startDate);
            if (endDate) params.append('end_date', endDate);


            const token = localStorage.getItem('jwtToken'); // Получаем токен из localStorage
            if (!token) {
                alert('Вы не авторизованы!');
                window.location.href = '/login'; // Перенаправляем на страницу логина
                return;
            }

            const response = await fetch(`/api/cars/filter`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) throw new Error('Ошибка при получении данных');
            const cars = await response.json();
            renderCars(cars);
        } catch (error) {
            console.error(error);
            alert('Не удалось загрузить автомобили');
        }
    }

    async function goToPresentationCarPage(id) {
        const token = localStorage.getItem("jwtToken");
        console.log('lololol')
        // Проверяем, авторизован ли пользователь
        if (!token) {
            alert("Вы не авторизованы!");
            window.location.href = "/login"; // Если не авторизованы, переходим на страницу логина
            return;
        }

        // const response = await fetch(`/api/car/page/${id}`, {
        //     method: 'GET',
        //     headers: {
        //         'Authorization': `Bearer ${token}`,
        //         'Content-Type': 'application/json'
        //     }
        // });

        // if (!response.ok) throw new Error('Ошибка при получении данных');
        // const cars = await response.json();

        // Перенаправляем пользователя на страницу поиска машин
        window.location.href = `/api/car/page/${id}`;
    }

    function renderCars(cars) {
        carList.innerHTML = '';

        if (cars.length === 0) {
            carList.innerHTML = '<p>Автомобили не найдены</p>';
            return;
        }

        cars.forEach(car => {
            const carItem = document.createElement('div');
            carItem.className = 'car-item';
            carItem.innerHTML = `
            <div id="goToPresentationCar">
                <img src="${car.image_url}" alt="${car.brand} ${car.model}" class="car-image">
                <p><strong>Бренд:</strong> ${car.brand}</p>
                <p><strong>Модель:</strong> ${car.model}</p>
                <p><strong>Год выпуска:</strong> ${car.year}</p>
                <p><strong>Цвет:</strong> ${car.color}</p>
                <p><strong>Пробег:</strong> ${car.mileage} км</p>
                <p><strong>Цена за день:</strong> ${car.price_per_day}₽</p>
                <p><strong>Статус:</strong> ${car.status}</p>
                <button class="action-button-car-presentation">перейти</button>
            </div>
        `;
            const button = carItem.querySelector('.action-button-car-presentation');
            button.addEventListener('click', () => {
                goToPresentationCarPage(car.id);
            });
            carList.appendChild(carItem);

        });
    }

    // Инициализация
    populateYears(filterYearFrom);
    populateYears(filterYearTo);
    setupDatePickers(startDateInput, endDateInput);
    loadBrands();

    filterBrand.addEventListener('change', () => {
        const selectedBrand = filterBrand.value;
        if (selectedBrand) {
            loadModels(selectedBrand);
        } else {
            filterModel.innerHTML = '<option value="">Модель</option>';
        }
    });

    const searchButton = document.getElementById('searchButton');
    searchButton.addEventListener('click', fetchCars);

    // Добавляем обработчик для кнопки профиля
    const profileIcon = document.getElementById('profileIcon');
    profileIcon.addEventListener('click', async () => {
        try {
            // Если запрос успешен, перенаправляем пользователя на страницу профиля
            window.location.href = '/api/profile/page';
        } catch (error) {
            alert("Ошибка при получении страницы профиля");
            console.error("Ошибка:", error);
        }
    });
});