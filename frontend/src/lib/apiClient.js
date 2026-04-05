/**
 * API Client for Parking Management System
 * Implements communication with backend REST API
 * Following Red/Green TDD methodology
 */

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

/**
 * Formats Date object to ISO 8601 string
 * @param {Date} date - Date to format
 * @returns {string} - ISO 8601 string
 * @throws {Error} - If date is invalid
 */
export function formatDateTimeForAPI(date) {
  if (!(date instanceof Date) || isNaN(date)) {
    throw new Error('Invalid date');
  }
  return date.toISOString();
}

/**
 * Parses ISO 8601 string to Date object
 * @param {string} isoString - ISO 8601 date string
 * @returns {Date} - Date object
 */
export function parseDateTimeFromAPI(isoString) {
  return new Date(isoString);
}

/**
 * Formats date for user display (DD.MM.YYYY HH:mm)
 * @param {Date|string} date - Date object or ISO string
 * @returns {string} - Formatted date string or "N/A" if invalid
 */
export function formatDateTimeForDisplay(date) {
  try {
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
  } catch {
    return 'N/A';
  }
}

/**
 * Validates phone number format (+XXXX...)
 * @param {string} phone - Phone number to validate
 * @returns {boolean} - True if valid
 */
export function validatePhoneNumber(phone) {
  if (!phone || typeof phone !== 'string') {
    return false;
  }
  // Can start with + and contain at least 9 digits
  const phoneRegex = /(\+)?\d{9,15}$/;
  return phoneRegex.test(phone);
}

/**
 * Validates license plate format
 * @param {string} plate - License plate to validate
 * @returns {boolean} - True if valid
 */
export function validateLicensePlate(plate) {
  if (!plate || typeof plate !== 'string') {
    return false;
  }
  // Remove spaces and check length (between 5-12 characters)
  const cleanPlate = plate.replace(/\s+/g, '');
  if (cleanPlate.length < 5 || cleanPlate.length > 12) {
    return false;
  }
  // Must contain letters and numbers
  return /^[A-Za-z0-9]+$/.test(cleanPlate);
}

/**
 * Validates client name
 * @param {string} name - Client name to validate
 * @returns {boolean} - True if valid
 */
export function validateClientName(name) {
  if (!name || typeof name !== 'string') {
    return false;
  }
  // Must have at least 1 non-space character and max 200 characters
  const trimmed = name.trim();
  return trimmed.length > 0 && trimmed.length <= 200;
}

/**
 * Validates parking spot number
 * @param {number} spot - Spot number to validate
 * @returns {boolean} - True if valid
 */
export function validateSpotNumber(spot) {
  if (typeof spot !== 'number') {
    return false;
  }
  // Must be positive integer
  return Number.isInteger(spot) && spot > 0;
}

/**
 * Creates a new parking session
 * @param {Object} sessionData - Session data
 * @param {string} sessionData.client_name - Client name
 * @param {string} sessionData.phone_number - Phone number
 * @param {string} sessionData.license_plate - License plate
 * @param {number} sessionData.spot_number - Parking spot number
 * @param {string} [sessionData.duration] - Optional session duration (e.g., "1h")
 * @returns {Promise<Object>} - API response
 * @throws {Error} - If request fails
 */
export async function createSession(sessionData) {
  const url = `${API_BASE_URL}/api/sessions`;
  const response = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(sessionData),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || `Failed to create session: ${response.statusText}`);
  }

  return data;
}

/**
 * Gets session information by phone number
 * @param {string} phoneNumber - Phone number
 * @returns {Promise<Object>} - Session data
 * @throws {Error} - If request fails or session not found
 */
export async function getSession(phoneNumber) {
  const encodedPhone = encodeURIComponent(phoneNumber);
  const url = `${API_BASE_URL}/api/sessions/${encodedPhone}`;
  const response = await fetch(url, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to get session: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Prolongs parking session by specified duration
 * @param {string} phoneNumber - Phone number
 * @param {string} duration - Duration to add (e.g., "1h")
 * @returns {Promise<Object>} - API response
 * @throws {Error} - If request fails
 */
export async function prolongSession(phoneNumber, duration) {
  const encodedPhone = encodeURIComponent(phoneNumber);
  const url = `${API_BASE_URL}/api/sessions/${encodedPhone}`;
  const response = await fetch(url, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ duration }),
  });

  if (!response.ok) {
    throw new Error(`Failed to prolong session: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Stops a parking session
 * @param {string} phoneNumber - Phone number
 * @returns {Promise<void>}
 * @throws {Error} - If request fails
 */
export async function stopSession(phoneNumber) {
  const encodedPhone = encodeURIComponent(phoneNumber);
  const url = `${API_BASE_URL}/api/sessions/${encodedPhone}`;
  const response = await fetch(url, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to stop session: ${response.statusText}`);
  }
}
