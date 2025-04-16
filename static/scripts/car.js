document.addEventListener('DOMContentLoaded', async () => {
    // Получаем ID автомобиля из URL
    const urlParams = new URLSearchParams(window.location.search);
    const carID = urlParams.get('id');

    // Проверяем наличие ID
    if (!carID) {
        handleError('Не указан ID автомобиля');
        return;
    }

    try {
        // Загружаем данные автомобиля с сервера
        const response = await fetch(`/api/cars/${carID}`);
        if (!response.ok) {
            throw new Error(`Ошибка HTTP: ${response.status}`);
        }

        const car = await response.json();

        // Заполняем страницу данными автомобиля
        renderCarDetails(car);
    } catch (error) {
        console.error(error);
        handleError('Не удалось загрузить данные автомобиля');
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

    carDetails.innerHTML = `
        <img src="${car.image_url}" alt="${car.brand} ${car.model}" class="car-image">
        <p><strong>Бренд:</strong> ${escapeHtml(car.brand)}</p>
        <p><strong>Модель:</strong> ${escapeHtml(car.model)}</p>
        <p><strong>Год выпуска:</strong> ${escapeHtml(car.year)}</p>
        <p><strong>Цвет:</strong> ${escapeHtml(car.color)}</p>
        <p><strong>Пробег:</strong> ${escapeHtml(car.mileage)} км</p>
        <p><strong>Цена за день:</strong> ${escapeHtml(car.price_per_day)}₽</p>
        <p><strong>Статус:</strong> ${escapeHtml(car.status)}</p>
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