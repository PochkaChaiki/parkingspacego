/**
 * API Service Tests - Red/Green TDD
 * Tests for parking management API client
 */
import {
  formatDateTimeForAPI,
  parseDateTimeFromAPI,
  formatDateTimeForDisplay,
  validatePhoneNumber,
  validateLicensePlate,
  validateClientName,
  validateSpotNumber,
  createSession,
  getSession,
  prolongSession,
  stopSession,
} from '../apiClient';

describe('Date and Time Formatting', () => {
  describe('formatDateTimeForAPI', () => {
    it('should format date to ISO 8601 format', () => {
      const date = new Date('2026-03-26T10:30:00Z');
      const result = formatDateTimeForAPI(date);
      expect(result).toBe('2026-03-26T10:30:00.000Z');
    });

    it('should throw error for null date', () => {
      expect(() => formatDateTimeForAPI(null)).toThrow('Invalid date');
    });

    it('should throw error for string date', () => {
      expect(() => formatDateTimeForAPI('2026-03-26')).toThrow('Invalid date');
    });

    it('should throw error for invalid date object', () => {
      expect(() => formatDateTimeForAPI(new Date('invalid'))).toThrow('Invalid date');
    });
  });

  describe('parseDateTimeFromAPI', () => {
    it('should parse ISO 8601 date string', () => {
      const isoString = '2026-03-26T10:30:00Z';
      const result = parseDateTimeFromAPI(isoString);
      expect(result).toBeInstanceOf(Date);
      expect(result.getFullYear()).toBe(2026);
      expect(result.getMonth()).toBe(2); // 0-indexed
      expect(result.getDate()).toBe(26);
    });

    it('should handle different ISO formats', () => {
      const result = parseDateTimeFromAPI('2026-03-26T10:30:00.000Z');
      expect(result).toBeInstanceOf(Date);
    });
  });

  describe('formatDateTimeForDisplay', () => {
    it('should format Date object to DD.MM.YYYY HH:mm', () => {
      const date = new Date('2026-03-26T10:30:00Z');
      const result = formatDateTimeForDisplay(date);
      expect(result).toMatch(/\d{2}\.\d{2}\.\d{4} \d{2}:\d{2}/);
    });

    it('should accept ISO string and format it', () => {
      const result = formatDateTimeForDisplay('2026-03-26T10:30:00Z');
      expect(result).toMatch(/\d{2}\.\d{2}\.\d{4} \d{2}:\d{2}/);
    });

    it('should return "N/A" for invalid date', () => {
      const result = formatDateTimeForDisplay('invalid');
      expect(result).toBe('N/A');
    });

    it('should return "N/A" for null', () => {
      const result = formatDateTimeForDisplay(null);
      expect(result).toBe('N/A');
    });
  });
});

describe('Validation Functions', () => {
  describe('validatePhoneNumber', () => {
    it('should validate correct phone number format', () => {
      expect(validatePhoneNumber('+79991234567')).toBe(true);
      expect(validatePhoneNumber('+1234567890')).toBe(true);
    });

    // it('should reject phone without plus sign', () => {
    //   expect(validatePhoneNumber('79991234567')).toBe(true);
    // });

    it('should reject empty phone', () => {
      expect(validatePhoneNumber('')).toBe(false);
      expect(validatePhoneNumber(null)).toBe(false);
    });

    it('should reject phone with letters', () => {
      expect(validatePhoneNumber('+7999ABC1234')).toBe(false);
    });

    it('should reject too short phone', () => {
      expect(validatePhoneNumber('+123')).toBe(false);
    });
  });

  describe('validateLicensePlate', () => {
    it('should validate correct license plate format', () => {
      expect(validateLicensePlate('A123BC140')).toBe(true);
      expect(validateLicensePlate('X999YZ77')).toBe(true);
    });

    it('should accept various formats', () => {
      expect(validateLicensePlate('AB123CD')).toBe(true);
      expect(validateLicensePlate('A 123 BC 140')).toBe(true);
    });

    it('should reject empty plate', () => {
      expect(validateLicensePlate('')).toBe(false);
      expect(validateLicensePlate(null)).toBe(false);
    });

    it('should reject too short plate', () => {
      expect(validateLicensePlate('A1')).toBe(false);
    });

    it('should reject too long plate', () => {
      expect(validateLicensePlate('ABCDEFGHIJKLM')).toBe(false); // 13 characters
    });
  });

  describe('validateClientName', () => {
    it('should validate non-empty name', () => {
      expect(validateClientName('Иван')).toBe(true);
      expect(validateClientName('John Doe')).toBe(true);
    });

    it('should reject empty name', () => {
      expect(validateClientName('')).toBe(false);
      expect(validateClientName(null)).toBe(false);
    });

    it('should reject name with only spaces', () => {
      expect(validateClientName('   ')).toBe(false);
    });

    it('should reject very long name', () => {
      expect(validateClientName('A'.repeat(201))).toBe(false);
    });
  });

  describe('validateSpotNumber', () => {
    it('should validate positive spot number', () => {
      expect(validateSpotNumber(1)).toBe(true);
      expect(validateSpotNumber(100)).toBe(true);
    });

    it('should reject zero', () => {
      expect(validateSpotNumber(0)).toBe(false);
    });

    it('should reject negative number', () => {
      expect(validateSpotNumber(-1)).toBe(false);
    });

    it('should reject non-integer', () => {
      expect(validateSpotNumber(1.5)).toBe(false);
    });

    it('should reject string', () => {
      expect(validateSpotNumber('42')).toBe(false);
    });
  });
});

describe('API Client Functions', () => {
  // Mock fetch
  beforeEach(() => {
    global.fetch = jest.fn();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('createSession', () => {
    it('should send POST request with session data', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({ status: 'success' }),
      });

      const sessionData = {
        client_name: 'Иван',
        phone_number: '79991234567',
        license_plate: 'A123BC140',
        spot_number: 42,
      };

      const result = await createSession(sessionData);
      expect(result).toEqual({ status: 'success' });
      expect(global.fetch).toHaveBeenCalled();
    });

    it('should construct correct API endpoint', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({ status: 'success' }),
      });

      const sessionData = {
        client_name: 'Иван',
        phone_number: '79991234567',
        license_plate: 'A123BC140',
        spot_number: 42,
      };

      await createSession(sessionData);
      const [url] = global.fetch.mock.calls[0];
      expect(url).toContain('/api/sessions');
    });

    it('should send POST method', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({ status: 'success' }),
      });

      const sessionData = {
        client_name: 'Иван',
        phone_number: '79991234567',
        license_plate: 'A123BC140',
        spot_number: 42,
      };

      await createSession(sessionData);
      const [, options] = global.fetch.mock.calls[0];
      expect(options.method).toBe('POST');
    });

    it('should throw error on network failure', async () => {
      global.fetch.mockRejectedValueOnce(new Error('Network error'));

      const sessionData = {
        client_name: 'Иван',
        phone_number: '79991234567',
        license_plate: 'A123BC140',
        spot_number: 42,
      };

      await expect(createSession(sessionData)).rejects.toThrow('Network error');
    });
  });

  describe('getSession', () => {
    it('should send GET request with phone number', async () => {
      const mockSession = {
        client_name: 'Иван',
        phone_number: '+79991234567',
        license_plate: 'A123BC140',
        spot_number: 42,
        start_time: '2026-03-26T10:30:00Z',
        end_time: '2026-03-26T11:30:00Z',
      };

      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockSession,
      });

      const result = await getSession('79991234567');
      expect(result).toEqual(mockSession);
    });

    it('should construct correct endpoint with phone number', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({}),
      });

      await getSession('79991234567');
      const [url] = global.fetch.mock.calls[0];
      expect(url).toContain('/api/sessions');
      expect(url).toContain(encodeURIComponent('79991234567')); // Phone is URL encoded
    });

    it('should throw error on 404', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => ({ error: 'Not found' }),
      });

      await expect(getSession('79991234567')).rejects.toThrow();
    });
  });

  describe('prolongSession', () => {
    it('should send PATCH request to prolong session', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ status: 'success' }),
      });

      const result = await prolongSession('79991234567', '1h');
      expect(result).toEqual({ status: 'success' });
    });

    it('should include duration in request body', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ status: 'success' }),
      });

      await prolongSession('79991234567', '2h');
      const [, options] = global.fetch.mock.calls[0];
      const body = JSON.parse(options.body);
      expect(body.duration).toBe('2h');
    });

    it('should use PATCH method', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ status: 'success' }),
      });

      await prolongSession('79991234567', '1h');
      const [, options] = global.fetch.mock.calls[0];
      expect(options.method).toBe('PATCH');
    });
  });

  describe('stopSession', () => {
    it('should send DELETE request to stop session', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 204,
      });

      await stopSession('79991234567');
      expect(global.fetch).toHaveBeenCalled();
    });

    it('should use DELETE method', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 204,
      });

      await stopSession('+79991234567');
      const [, options] = global.fetch.mock.calls[0];
      expect(options.method).toBe('DELETE');
    });

    it('should construct correct endpoint', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        status: 204,
      });

      await stopSession('79991234567');
      const [url] = global.fetch.mock.calls[0];
      expect(url).toContain('/api/sessions');
      expect(url).toContain(encodeURIComponent('79991234567')); // Phone is URL encoded
    });

    it('should throw error on failure', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
      });

      await expect(stopSession('79991234567')).rejects.toThrow();
    });
  });
});
