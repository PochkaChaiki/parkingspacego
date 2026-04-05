/**
 * Основное приложение для работы с парковкой
 * Связывает DOM события с логикой API
 */

const API_BASE_URL = process.env.API_URL || 'http://localhost:8080';

// ===== Inline API Functions =====

/**
 * Форматирует дату и время для отправки на сервер
 */
function formatDateTimeForAPI(date) {
    if (!(date instanceof Date) || isNaN(date)) {
        throw new Error('Invalid date');
    }
    return date.toISOString();
}

/**
 * Форматирует дату для отображения пользователю
 */
function formatDateTimeForDisplay(date) {
    if (typeof date === 'string') {
        date = new Date(date);
    }
    if (!(date instanceof Date) || isNaN(date)) {
        return 'N/A';
    }
    const pad = (n) => String(n).padStart(2, '0');
    const year = date.getFullYear();
    const month = pad(date.getMonth() + 1);
    const day = pad(date.getDate());
    const hours = pad(date.getHours());
    const minutes = pad(date.getMinutes());
    return `${day}.${month}.${year} ${hours}:${minutes}`;
}

/**
 * Валидирует номер телефона
 */
function validatePhoneNumber(phone) {
    const phoneRegex = /^\+\d{1,15}$/;
    return phoneRegex.test(phone);
}

/**
 * Валидирует номер автомобиля
 */
function validateCarNumber(carNumber) {
    return carNumber.length > 0 && carNumber.length <= 20;
}

/**
 * Валидирует номер парковочного места
 */
function validateSpotNumber(spotNumber) {
    const num = parseInt(spotNumber, 10);
    return !isNaN(num) && num > 0 && num <= 9999;
}

/**
 * Создает новую парковочную резервацию
 */
async function createRecord(record) {
    const response = await fetch(`${API_BASE_URL}/api/records`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(record),
    });

    if (!response.ok) {
        const error = await response.text();
        throw new Error(`Failed to create record: ${error}`);
    }

    return response.json();
}

/**
 * Получает все парковочные резервации
 */
async function getAllRecords() {
    const response = await fetch(`${API_BASE_URL}/api/records`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });

    if (!response.ok) {
        throw new Error('Failed to fetch records');
    }

    return response.json();
}

/**
 * Завершает парковочную резервацию
 */
async function completeRecord(recordId, endTime) {
    const response = await fetch(`${API_BASE_URL}/api/records/${recordId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            end_time: formatDateTimeForAPI(endTime),
            status: 'completed',
        }),
    });

    if (!response.ok) {
        throw new Error('Failed to complete record');
    }

    return response.json();
}

/**
 * Удаляет все резервации по номеру телефона
 */
async function deleteRecordsByPhone(phoneNumber) {
    const response = await fetch(`${API_BASE_URL}/api/records/${phoneNumber}`, {
        method: 'DELETE',
    });

    if (!response.ok) {
        throw new Error('Failed to delete records');
    }
}

/**
 * Получает резервации по номеру телефона
 */
async function getRecordsByPhone(phoneNumber) {
    const allRecords = await getAllRecords();
    return allRecords.filter(
        (record) => record.phone_number === phoneNumber
    );
}

// ===== DOM Application Code =====

// DOM элементы
const bookingForm = document.getElementById('bookingForm');
const formMessage = document.getElementById('formMessage');
const recordsList = document.getElementById('recordsList');
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
function createRecordCard(record) {
    const statusClass = `status-${record.status}`;
    const endTimeDisplay = record.end_time
        ? formatDateTimeForDisplay(record.end_time)
        : '—';

    const card = document.createElement('div');
    card.className = 'record-card';
    card.innerHTML = `
        <div class="record-header">
            <div class="record-title">${escapeHtml(record.car_number)}</div>
            <span class="record-status ${statusClass}">${record.status}</span>
        </div>
        <div class="record-info">
            <div class="info-row">
                <span class="info-label">Клиент:</span>
                <span class="info-value">${escapeHtml(record.client_name)}</span>
            </div>
            <div class="info-row">
                <span class="info-label">Телефон:</span>
                <span class="info-value">${escapeHtml(record.phone_number)}</span>
            </div>
            <div class="info-row">
                <span class="info-label">Место:</span>
                <span class="info-value">#${record.spot_number}</span>
            </div>
            <div class="info-row">
                <span class="info-label">Начало:</span>
                <span class="info-value">${formatDateTimeForDisplay(record.start_time)}</span>
            </div>
            <div class="info-row">
                <span class="info-label">Окончание:</span>
                <span class="info-value">${endTimeDisplay}</span>
            </div>
        </div>
        <div class="record-actions">
            ${
                record.status === 'active'
                    ? `<button class="btn btn-success complete-btn" data-id="${record.id}">Завершить</button>`
                    : ''
            }
            <button class="btn btn-danger delete-btn" data-phone="${record.phone_number}">Удалить</button>
        </div>
    `;

    // Обработчики событий
    const completeBtn = card.querySelector('.complete-btn');
    if (completeBtn) {
        completeBtn.addEventListener('click', async (e) => {
            e.preventDefault();
            await handleCompleteRecord(record.id);
        });
    }

    const deleteBtn = card.querySelector('.delete-btn');
    if (deleteBtn) {
        deleteBtn.addEventListener('click', async (e) => {
            e.preventDefault();
            await handleDeleteByPhone(record.phone_number);
        });
    }

    return card;
}

/**
 * Загружает и отображает все резервации
 */
async function loadAllRecords() {
    try {
        recordsList.innerHTML = '<div class="loading"><div class="spinner"></div></div>';
        const records = await getAllRecords();

        if (!records || records.length === 0) {
            recordsList.innerHTML =
                '<div class="empty-state"><p>Нет активных резервацій</p></div>';
            return;
        }

        recordsList.innerHTML = '';
        records.forEach((record) => {
            recordsList.appendChild(createRecordCard(record));
        });
    } catch (error) {
        console.error('Error loading records:', error);
        recordsList.innerHTML =
            '<div class="empty-state"><p>Ошибка при загрузке резервацій</p></div>';
    }
}

/**
 * Обработчик отправки формы создания резервации
 */
bookingForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const clientName = document.getElementById('clientName').value.trim();
    const phoneNumber = document.getElementById('phoneNumber').value.trim();
    const carNumber = document.getElementById('carNumber').value.trim();
    const spotNumber = parseInt(document.getElementById('spotNumber').value, 10);
    const startTimeInput = document.getElementById('startTime').value;

    // Валидация
    if (!clientName) {
        showMessage('Введите имя клиента', 'error');
        return;
    }

    if (!validatePhoneNumber(phoneNumber)) {
        showMessage('Введите корректный номер телефона (+XXXXXXXXXXX)', 'error');
        return;
    }

    if (!validateCarNumber(carNumber)) {
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

        const record = await createRecord({
            client_name: clientName,
            phone_number: phoneNumber,
            car_number: carNumber,
            spot_number: spotNumber,
            start_time: formatDateTimeForAPI(startTime),
        });

        showMessage(`✓ Резервация создана (ID: ${record.id})`);
        bookingForm.reset();
        await loadAllRecords();
    } catch (error) {
        console.error('Error creating record:', error);
        showMessage(`✗ Ошибка: ${error.message}`, 'error');
    }
});

/**
 * Обработчик завершения парковки
 */
async function handleCompleteRecord(recordId) {
    try {
        const endTime = new Date();
        const updated = await completeRecord(recordId, endTime);
        showMessage('✓ Парковка завершена');
        await loadAllRecords();
    } catch (error) {
        console.error('Error completing record:', error);
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
        await deleteRecordsByPhone(phoneNumber);
        showMessage('✓ Резервации удалены');
        await loadAllRecords();
    } catch (error) {
        console.error('Error deleting records:', error);
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
        const records = await getRecordsByPhone(searchPhone);

        if (records.length === 0) {
            searchResults.innerHTML =
                '<div class="empty-state"><p>Резервации не найдены</p></div>';
            return;
        }

        searchResults.innerHTML = '';
        records.forEach((record) => {
            searchResults.appendChild(createRecordCard(record));
        });
    } catch (error) {
        console.error('Error searching records:', error);
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
    loadAllRecords();

    // Обновляем список каждые 30 секунд
    setInterval(loadAllRecords, 30000);
});
