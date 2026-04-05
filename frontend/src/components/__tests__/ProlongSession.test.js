/**
 * Component Tests - Prolong Session
 * Following Red/Green TDD methodology
 */
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import ProlongSession from '../ProlongSession';
import * as apiClient from '@/lib/apiClient';

jest.mock('@/lib/apiClient');

describe('ProlongSession Component', () => {
  const mockOnSuccess = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render form with phone and duration fields', () => {
    render(<ProlongSession onSuccess={mockOnSuccess} />);

    expect(screen.getByText(/Prolong Parking Session/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Phone Number/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Duration to Add/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Prolong Session/i })).toBeInTheDocument();
  });

  it('should have default duration value of 1h', () => {
    render(<ProlongSession onSuccess={mockOnSuccess} />);
    const durationInput = screen.getByLabelText(/Duration to Add/i);
    expect(durationInput).toHaveValue('1h');
  });

  it('should disable submit button when phone number is empty', () => {
    render(<ProlongSession onSuccess={mockOnSuccess} />);
    const submitButton = screen.getByRole('button', { name: /Prolong Session/i });
    expect(submitButton).toBeDisabled();
  });

  it('should enable submit button when phone number is filled', async () => {
    render(<ProlongSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const submitButton = screen.getByRole('button', { name: /Prolong Session/i });

    await waitFor(() => {
      expect(submitButton).not.toBeDisabled();
    });
  });

  it('should call prolongSession API when form is submitted', async () => {
    apiClient.prolongSession.mockResolvedValueOnce({ status: 'success' });

    render(<ProlongSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const submitButton = screen.getByRole('button', { name: /Prolong Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(apiClient.prolongSession).toHaveBeenCalledWith('+79991234567', '1h');
    });
  });

  it('should use custom duration when provided', async () => {
    apiClient.prolongSession.mockResolvedValueOnce({ status: 'success' });

    render(<ProlongSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const durationInput = screen.getByLabelText(/Duration to Add/i);
    fireEvent.change(durationInput, { target: { value: '2h' } });

    const submitButton = screen.getByRole('button', { name: /Prolong Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(apiClient.prolongSession).toHaveBeenCalledWith('+79991234567', '2h');
    });
  });

  it('should call onSuccess callback on successful submission', async () => {
    apiClient.prolongSession.mockResolvedValueOnce({ status: 'success' });

    render(<ProlongSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const submitButton = screen.getByRole('button', { name: /Prolong Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalled();
    });
  });

  it('should show error message on submission failure', async () => {
    const errorMessage = 'Session not found';
    apiClient.prolongSession.mockRejectedValueOnce(new Error(errorMessage));

    render(<ProlongSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const submitButton = screen.getByRole('button', { name: /Prolong Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(new RegExp(errorMessage, 'i'))).toBeInTheDocument();
    });
  });

  it('should reset form after successful submission', async () => {
    apiClient.prolongSession.mockResolvedValueOnce({ status: 'success' });

    render(<ProlongSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const submitButton = screen.getByRole('button', { name: /Prolong Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(phoneInput.value).toBe('');
    });
  });

  it('should display loading state during submission', async () => {
    apiClient.prolongSession.mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve({ status: 'success' }), 100))
    );

    render(<ProlongSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const submitButton = screen.getByRole('button', { name: /Prolong Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Prolonging/i })).toBeInTheDocument();
    });
  });
});
