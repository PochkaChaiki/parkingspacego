/**
 * Основное приложение для работы с парковкой
 * Связывает DOM события с логикой API
 */

// Import logic functions (ES modules)
import {
    formatDateTimeForAPI,
    formatDateTimeForDisplay,
    validatePhoneNumber,
    validateLicensePlate,
    validateSpotNumber,
    createSession,
    getAllSessions,
    completeSession,
    deleteSessionsByPhone,
    getSessionsByPhone,
} from './logic/api.js';

// DOM элементы
const bookingForm = document.getElementById('bookingForm');
const formMessage = document.getElementById('formMessage');
const sessionsList = document.getElementById('sessionsList');
const searchForm = document.getElementById('searchForm');
const searchResults = document.getElementById('searchResults');

/**
 * Показывает сообщение в форме
 */
function showMessage(message, type = 'success') {
    formMessage.textContent = message;
    formMessage.className = `message ${type}`;
    if (type === 'success') {
        setTimeout(() => {
            formMessage.className = 'message';
        }, 3000);
    }
}

/**
 * Создает HTML карточку для одной резервации
 */
function createSessionCard(session) {
    const statusClass = `status-${session.status}`;
    const endTimeDisplay = session.end_time
        ? formatDateTimeForDisplay(session.end_time)
        : '—';

    const card = document.createElement('div');
    card.className = 'session-card';
    card.innerHTML = `
        <div class="session-header">
            <div class="session-title">${escapeHtml(session.license_plate)}</div>
            <span class="session-status ${statusClass}">${session.status}</span>
        </div>
        <div class="session-info">
            <div class="info-row">
                <span class="info-label">Клиент:</span>
                <span class="info-value">${escapeHtml(session.client_name)}</span>
            </div>
            <div class="info-row">
                <span class="info-label">Телефон:</span>
                <span class="info-value">${escapeHtml(session.phone_number)}</span>
            </div>
            <div class="info-row">
                <span class="info-label">Место:</span>
                <span class="info-value">#${session.spot_number}</span>
            </div>
            <div class="info-row">
                <span class="info-label">Начало:</span>
                <span class="info-value">${formatDateTimeForDisplay(session.start_time)}</span>
            </div>
            <div class="info-row">
                <span class="info-label">Окончание:</span>
                <span class="info-value">${endTimeDisplay}</span>
            </div>
        </div>
        <div class="session-actions">
            ${
                session.status === 'active'
                    ? `<button class="btn btn-success complete-btn" data-id="${session.id}">Завершить</button>`
                    : ''
            }
            <button class="btn btn-danger delete-btn" data-phone="${session.phone_number}">Удалить</button>
        </div>
    `;

    // Обработчики событий
    const completeBtn = card.querySelector('.complete-btn');
    if (completeBtn) {
        completeBtn.addEventListener('click', async (e) => {
            e.preventDefault();
            await handleCompleteSymbol(session.id);
        });
    }

    const deleteBtn = card.querySelector('.delete-btn');
    if (deleteBtn) {
        deleteBtn.addEventListener('click', async (e) => {
            e.preventDefault();
            await handleDeleteByPhone(session.phone_number);
        });
    }

    return card;
}

/**
 * Загружает и отображает все резервации
 */
async function loadAllSessions() {
    try {
        sessionsList.innerHTML = '<div class="loading"><div class="spinner"></div></div>';
        const sessions = await getAllSessions();

        if (!sessions || sessions.length === 0) {
            sessionsList.innerHTML =
                '<div class="empty-state"><p>Нет активных резерваций</p></div>';
            return;
        }

        sessionsList.innerHTML = '';
        sessions.forEach((session) => {
            sessionsList.appendChild(createSessionCard(session));
        });
    } catch (error) {
        console.error('Error loading sessions:', error);
        sessionsList.innerHTML =
            '<div class="empty-state"><p>Ошибка при загрузке резерваций</p></div>';
    }
}

/**
 * Обработчик отправки формы создания резервации
 */
bookingForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const clientName = document.getElementById('clientName').value.trim();
    const phoneNumber = document.getElementById('phoneNumber').value.trim();
    const licensePlate = document.getElementById('licensePlate').value.trim();
    const spotNumber = parseInt(document.getElementById('spotNumber').value, 10);
    const startTimeInput = document.getElementById('startTime').value;

    // Валидация
    if (!clientName) {
        showMessage('Введите имя клиента', 'error');
        return;
    }

    if (!validatePhoneNumber(phoneNumber)) {
        showMessage('Введите корректный номер телефона (XXXXXXXXXXX)', 'error');
        return;
    }

    if (!validateLicensePlate(licensePlate)) {
        showMessage('Введите номер автомобиля', 'error');
        return;
    }

    if (!validateSpotNumber(spotNumber)) {
        showMessage('Введите валидный номер места (1-9999)', 'error');
        return;
    }

    if (!startTimeInput) {
        showMessage('Выберите время начала', 'error');
        return;
    }

    try {
        const startTime = new Date(startTimeInput);
        if (isNaN(startTime)) {
            showMessage('Ошибка форматирования времени', 'error');
            return;
        }

        const session = await createSession({
            client_name: clientName,
            phone_number: phoneNumber.replaceAll(" ", ""),
            license_plate: licensePlate,
            spot_number: spotNumber,
            start_time: formatDateTimeForAPI(startTime),
        });

        showMessage(`✓ Резервация создана (ID: ${session.id})`);
        bookingForm.reset();
        await loadAllSessions();
    } catch (error) {
        console.error('Error creating session:', error);
        showMessage(`✗ Ошибка: ${error.message}`, 'error');
    }
});

/**
 * Обработчик завершения парковки
 */
async function handleCompleteSymbol(sessionId) {
    try {
        const endTime = new Date();
        const updated = await completeSession(sessionId, endTime);
        showMessage('✓ Парковка завершена');
        await loadAllSessions();
    } catch (error) {
        console.error('Error completing session:', error);
        showMessage(`✗ Ошибка: ${error.message}`, 'error');
    }
}

/**
 * Обработчик удаления резерваций по телефону
 */
async function handleDeleteByPhone(phoneNumber) {
    if (!confirm(`Удалить все резервации для ${phoneNumber}?`)) {
        return;
    }

    try {
        await deleteSessionsByPhone(phoneNumber);
        showMessage('✓ Резервации удалены');
        await loadAllSessions();
    } catch (error) {
        console.error('Error deleting session:', error);
        showMessage(`✗ Ошибка: ${error.message}`, 'error');
    }
}

/**
 * Обработчик формы поиска
 */
searchForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const searchPhone = document.getElementById('searchPhone').value.trim();

    if (!validatePhoneNumber(searchPhone)) {
        showMessage('Введите корректный номер телефона', 'error');
        return;
    }

    try {
        searchResults.innerHTML = '<div class="loading"><div class="spinner"></div></div>';
        const sessions = await getSessionsByPhone(searchPhone.replaceAll(" ", ""));

        if (sessions.length === 0) {
            searchResults.innerHTML =
                '<div class="empty-state"><p>Резервации не найдены</p></div>';
            return;
        }

        searchResults.innerHTML = '';
        sessions.forEach((session) => {
            searchResults.appendChild(createSessionCard(session));
        });
    } catch (error) {
        console.error('Error searching sessions:', error);
        searchResults.innerHTML =
            '<div class="empty-state"><p>Ошибка при поиске</p></div>';
    }
});

/**
 * Экранирует HTML для безопасного вывода
 */
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

/**
 * Инициализирует приложение при загрузке страницы
 */
document.addEventListener('DOMContentLoaded', () => {
    // Устанавливаем текущее время по умолчанию
    const now = new Date();
    const isoDateTime = now.toISOString().slice(0, 16);
    document.getElementById('startTime').value = isoDateTime;

    // Загружаем все резервации
    loadAllSessions();

    // Обновляем список каждые 30 секунд
    setInterval(loadAllSessions, 30000);
});
