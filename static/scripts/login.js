document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    try {
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        const response = await fetch('/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        });

        if (response.ok) {
            const data = await response.json();
            console.log("Токен:", data.token); // Проверяем, получаем ли токен

            localStorage.setItem('jwtToken', data.token);
            console.log("Перенаправление...");
            window.location.href = '/api/profile/page';
        } else {
            console.log("Ошибка авторизации");
            document.getElementById('errorMessage').style.display = 'block';
        }
    } catch (error) {
        console.error('Ошибка при входе:', error);
    }
});