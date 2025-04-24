document.addEventListener('DOMContentLoaded', async () => {
    // Получаем JWT-токен из localStorage
    const token = localStorage.getItem("jwtToken");
    // Проверяем авторизацию
    if (!token) {
        alert("Вы не авторизованы!");
        window.location.href = "/login"; // Переход на страницу логина, если токен не найден
        return;
    }

    // Получаем ID автомобиля из URL
    const pathParts = window.location.pathname.split('/'); // Разбиваем путь на части
    const carID = pathParts[pathParts.length - 1]; // Берём последнюю часть пути
    console.log('ID автомобиля:', carID);

    // Проверяем наличие ID
    if (!carID || isNaN(carID)) {
        handleError('Не указан или некорректен ID автомобиля');
        return;
    }

    try {
        // Отправляем запрос на бэкенд для получения данных о машине
        const response = await fetch(`/api/cars/${carID}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}` // Добавляем токен в заголовок
            }
        });

        // Проверяем статус ответа
        if (response.status === 404) {
            throw new Error("Машина не найдена");
        }
        if (!response.ok) {
            throw new Error(`Ошибка ${response.status}: ${await response.text()}`);
        }

        // Парсим данные автомобиля
        const car = await response.json();

        // Заполняем страницу данными автомобиля
        renderCarDetails(car);

        // Инициализация кнопки "Оформить заказ"
        initRentModal(car, carID, token);

        // Инициализация отзывов
        initReviews(carID, token);

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
    } catch (error) {
        alert("Ошибка загрузки данных автомобиля");
        console.error("Ошибка:", error);
    }
});

/**
 * Обрабатывает ошибки и показывает сообщение пользователю.
 * @param {string} message - Сообщение об ошибке.
 */
function handleError(message) {
    alert(message); // Можно заменить на более красивое уведомление
    window.location.href = '/'; // Перенаправляем на главную страницу
}

/**
 * Заполняет страницу данными автомобиля.
 * @param {Object} car - Данные автомобиля.
 */
function renderCarDetails(car) {
    const carDetails = document.getElementById('carDetails');
    if (!carDetails) {
        console.error('Элемент carDetails не найден');
        return;
    }

    // Форматируем даты для удобства чтения
    const createdAt = formatDate(car.created_at);
    const updatedAt = formatDate(car.updated_at);

    carDetails.innerHTML = `
        <img src="${car.image_url}" alt="${car.brand} ${car.model}" class="car-image">
        <p><strong>Бренд:</strong> ${escapeHtml(car.brand)}</p>
        <p><strong>Модель:</strong> ${escapeHtml(car.model)}</p>
        <p><strong>Год выпуска:</strong> ${escapeHtml(car.year)}</p>
        <p><strong>Цвет:</strong> ${escapeHtml(car.color)}</p>
        <p><strong>Пробег:</strong> ${escapeHtml(car.mileage)} км</p>
        <p><strong>Цена за день:</strong> ${escapeHtml(car.price_per_day)}₽</p>
        <p><strong>Статус:</strong> ${escapeHtml(car.status)}</p>
        <p><strong>ID локации:</strong> ${escapeHtml(car.location_id)}</p>
        <p><strong>Дата создания:</strong> ${createdAt}</p>
        <p><strong>Дата обновления:</strong> ${updatedAt}</p>
    `;
}

/**
 * Экранирует HTML-символы для предотвращения XSS-атак.
 * @param {string} str - Входная строка.
 * @returns {string} - Безопасная строка.
 */
function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

/**
 * Форматирует дату для удобства чтения.
 * @param {string} dateString - Дата в формате ISO 8601.
 * @returns {string} - Отформатированная дата.
 */
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString(); // Форматирует дату под локаль пользователя
}

/**
 * Инициализирует модальное окно для оформления заказа.
 * @param {Object} car - Данные автомобиля.
 * @param {string} carID - ID автомобиля.
 * @param {string} token - JWT-токен.
 */
function initRentModal(car, carID, token) {
    // Найдите существующую кнопку
    const rentButton = document.getElementById('rentButton');

    // Модальное окно
    const rentModal = document.createElement('div');
    rentModal.id = 'rentModal';
    rentModal.className = 'modal';
    rentModal.style.display = 'none';
    rentModal.innerHTML = `
        <div class="modal-content">
            <span class="close">&times;</span>
            <h2>Оформление заказа</h2>
            <form id="rentForm">
                <label for="startDate">Дата начала аренды:</label>
                <input type="date" id="startDate" name="startDate" required>
                <label for="endDate">Дата окончания аренды:</label>
                <input type="date" id="endDate" name="endDate" required>
                <p><strong>Итоговая стоимость:</strong> <span id="totalCost">0 ₽</span></p>
                <button type="submit" id="confirmRent" class="confirm-rent-button">Подтвердить заказ</button>
            </form>
        </div>
    `;

    // Добавляем модальное окно в DOM
    document.body.appendChild(rentModal);

    // Ссылки на элементы модального окна
    const closeBtn = rentModal.querySelector('.close');
    const startDateInput = rentModal.querySelector('#startDate');
    const endDateInput = rentModal.querySelector('#endDate');
    const totalCostSpan = rentModal.querySelector('#totalCost');
    const rentForm = rentModal.querySelector('#rentForm');

    // Открываем модальное окно
    rentButton.addEventListener('click', () => {
        rentModal.style.display = 'flex';
    });

    // Закрываем модальное окно
    closeBtn.addEventListener('click', () => {
        rentModal.style.display = 'none';
    });

    // Закрываем модальное окно при клике вне его области
    window.addEventListener('click', (event) => {
        if (event.target === rentModal) {
            rentModal.style.display = 'none';
        }
    });

    // Расчет стоимости
    function calculateTotalCost(startDate, endDate) {
        const days = calculateDays(startDate, endDate);
        return days * car.price_per_day;
    }

    function calculateDays(start, end) {
        const startMs = new Date(start).getTime();
        const endMs = new Date(end).getTime();
        return Math.ceil((endMs - startMs) / (1000 * 60 * 60 * 24));
    }

    startDateInput.addEventListener('change', validateDates);
    endDateInput.addEventListener('change', validateDates);

    function validateDates() {
        const startDate = startDateInput.value;
        const endDate = endDateInput.value;
        if (startDate && endDate) {
            const totalCost = calculateTotalCost(startDate, endDate);
            totalCostSpan.textContent = `${totalCost} ₽`;
        }
    }

    // Подтверждение заказа
    rentForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const startDate = startDateInput.value;
        const endDate = endDateInput.value;
        const totalCost = parseInt(totalCostSpan.textContent);
        if (!startDate || !endDate || isNaN(totalCost)) {
            alert("Пожалуйста, выберите корректные даты.");
            return;
        }

        try {
            const response = await fetch('/api/orders', {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
                body: JSON.stringify({
                    car_id: parseInt(carID, 10),
                    start_date: startDate,
                    end_date: endDate,
                    total_cost: totalCost
                })
            });

            if (!response.ok) {
                throw new Error(`Ошибка ${response.status}: ${await response.text()}`);
            }

            alert("Заказ успешно оформлен!");
            rentModal.style.display = 'none';
        } catch (error) {
            alert("Ошибка при оформлении заказа");
            console.error("Ошибка:", error);
        }
    });
}

/**
 * Инициализирует функционал отзывов.
 * @param {string} carID - ID автомобиля.
 * @param {string} token - JWT-токен.
 */
function initReviews(carID, token) {
    // Загружаем отзывы
    loadReviews(carID);

    // Находим форму для отправки отзыва
    const reviewForm = document.getElementById('addReviewForm');

    // Предполагаем, что carID уже известен (например, из URL или глобальной переменной)
    if (!carID) {
        alert("Car ID не найден. Пожалуйста, вернитесь на страницу автомобиля.");
    }

    if (reviewForm) {
        reviewForm.addEventListener('submit', async (e) => {
            e.preventDefault(); // Предотвращаем стандартное поведение формы

            // Собираем данные из формы
            const rating = parseInt(document.getElementById('rating').value.trim(), 10);
            const comment = document.getElementById('comment').value.trim();

            // Проверяем, что все обязательные поля заполнены
            if (!rating || !comment) {
                alert("Пожалуйста, заполните все поля.");
                return;
            }

            // Проверяем, что рейтинг находится в допустимом диапазоне (например, от 1 до 5)
            if (rating < 1 || rating > 5) {
                alert("Рейтинг должен быть числом от 1 до 5.");
                return;
            }

            try {
                // Отправляем запрос на сервер
                const response = await fetch(`/api/reviews`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${token}` // Токен авторизации
                    },
                    body: JSON.stringify({
                        car_id: Math.abs(parseInt(carID, 10)), // Преобразуем carID в uint (беззнаковое целое)
                        rating: parseInt(rating, 10),          // Преобразуем rating в int (целое со знаком)
                        comment: comment                       // Комментарий остается строкой
                    })
                });

                // Проверяем статус ответа
                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Ошибка ${response.status}`);
                }

                // Очищаем форму после успешной отправки
                document.getElementById('rating').value = "";
                document.getElementById('comment').value = "";

                alert("Отзыв успешно добавлен!");
                loadReviews(carID); // Обновляем список отзывов
            } catch (error) {
                alert("Ошибка при отправке отзыва");
                console.error("Ошибка:", error);
            }
        });
    }
}

// Функция для загрузки отзывов
async function loadReviews(carID) {
    try {
        const token = localStorage.getItem("jwtToken");
        const response = await fetch(`/api/reviews/${carID}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`Ошибка ${response.status}: ${await response.text()}`);
        }

        const reviews = await response.json();
        const reviewsList = document.getElementById('reviewsList');
        reviewsList.innerHTML = ''; // Очищаем текущий список

        // Добавляем отзывы в DOM
        reviews.forEach(review => {
            const reviewElement = document.createElement('div');
            reviewElement.classList.add('review-item');
            reviewElement.innerHTML = `
                <p><strong>Рейтинг:</strong> ${review.rating}</p>
                <p><strong>Комментарий:</strong> ${review.comment}</p>
                <p><em>Добавлено: ${new Date(review.created_at).toLocaleString()}</em></p>
                <hr>
            `;
            reviewsList.appendChild(reviewElement);
        });
    } catch (error) {
        console.error("Ошибка при загрузке отзывов:", error);
    }
}