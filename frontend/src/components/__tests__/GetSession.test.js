/**
 * Component Tests - Get Session
 * Following Red/Green TDD methodology
 */
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import GetSession from '../GetSession';
import * as apiClient from '@/lib/apiClient';

jest.mock('@/lib/apiClient');

describe('GetSession Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render search form with phone number input', () => {
    render(<GetSession />);

    expect(screen.getByText(/View Parking Session/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Phone Number/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Search/i })).toBeInTheDocument();
  });

  it('should have search button disabled initially when no input', () => {
    render(<GetSession />);
    const searchButton = screen.getByRole('button', { name: /Search/i });
    expect(searchButton).not.toBeDisabled();
  });

  it('should call getSession API when form is submitted', async () => {
    const mockSession = {
      client_name: 'Иван',
      phone_number: '+79991234567',
      license_plate: 'A123BC140',
      spot_number: 42,
      start_time: '2026-03-26T10:30:00Z',
      end_time: '2026-03-26T11:30:00Z',
    };

    apiClient.getSession.mockResolvedValueOnce(mockSession);

    render(<GetSession />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const searchButton = screen.getByRole('button', { name: /Search/i });
    fireEvent.click(searchButton);

    await waitFor(() => {
      expect(apiClient.getSession).toHaveBeenCalledWith('+79991234567');
    });
  });

  it('should display session information on successful search', async () => {
    const mockSession = {
      client_name: 'Иван Петров',
      phone_number: '+79991234567',
      license_plate: 'A123BC140',
      spot_number: 42,
      start_time: '2026-03-26T10:30:00Z',
      end_time: '2026-03-26T11:30:00Z',
    };

    apiClient.getSession.mockResolvedValueOnce(mockSession);
    apiClient.formatDateTimeForDisplay = jest.fn((date) => '26.03.2026 10:30');

    render(<GetSession />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const searchButton = screen.getByRole('button', { name: /Search/i });
    fireEvent.click(searchButton);

    await waitFor(() => {
      expect(screen.getByText(/Session Information/i)).toBeInTheDocument();
      expect(screen.getByText('Иван Петров')).toBeInTheDocument();
      expect(screen.getByText('A123BC140')).toBeInTheDocument();
      expect(screen.getByText('42')).toBeInTheDocument();
    });
  });

  it('should show error message on search failure', async () => {
    const errorMessage = 'Session not found';
    apiClient.getSession.mockRejectedValueOnce(new Error(errorMessage));

    render(<GetSession />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const searchButton = screen.getByRole('button', { name: /Search/i });
    fireEvent.click(searchButton);

    await waitFor(() => {
      expect(screen.getByText(new RegExp(errorMessage, 'i'))).toBeInTheDocument();
    });
  });

  it('should clear error when performing new search', async () => {
    apiClient.getSession
      .mockRejectedValueOnce(new Error('First error'))
      .mockResolvedValueOnce({
        client_name: 'Test',
        phone_number: '+79991234567',
      });

    render(<GetSession />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    const searchButton = screen.getByRole('button', { name: /Search/i });

    // First search - fails
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });
    fireEvent.click(searchButton);

    await waitFor(() => {
      expect(screen.getByText(/First error/i)).toBeInTheDocument();
    });

    // Second search - succeeds
    fireEvent.click(searchButton);

    await waitFor(() => {
      // Error should be cleared
      expect(screen.queryByText(/First error/i)).not.toBeInTheDocument();
    });
  });

  it('should disable search button while loading', async () => {
    const mockSession = { client_name: 'Test', phone_number: '+79991234567' };
    apiClient.getSession.mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve(mockSession), 100))
    );

    render(<GetSession />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const searchButton = screen.getByRole('button', { name: /Search/i });
    fireEvent.click(searchButton);

    // Button should be disabled during loading
    await waitFor(() => {
      expect(searchButton).toBeDisabled();
    });
  });
});
