/**
 * Component Tests - Stop Session
 * Following Red/Green TDD methodology
 */
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import StopSession from '../StopSession';
import * as apiClient from '@/lib/apiClient';

jest.mock('@/lib/apiClient');

describe('StopSession Component', () => {
  const mockOnSuccess = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render form with phone number input', () => {
    render(<StopSession onSuccess={mockOnSuccess} />);

    expect(screen.getByText(/Stop Parking Session/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Phone Number/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Stop Session/i })).toBeInTheDocument();
  });

  it('should have stop button disabled initially', () => {
    render(<StopSession onSuccess={mockOnSuccess} />);
    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    expect(stopButton).toBeDisabled();
  });

  it('should enable stop button when phone number is provided', async () => {
    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });

    await waitFor(() => {
      expect(stopButton).not.toBeDisabled();
    });
  });

  it('should show confirmation dialog when stop button is clicked', async () => {
    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButton);

    await waitFor(() => {
      expect(screen.getByText(/Are you sure you want to stop the session/i)).toBeInTheDocument();
    });
  });

  it('should show confirmation with phone number', async () => {
    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButton);

    await waitFor(() => {
      expect(screen.getByText(/\+79991234567/)).toBeInTheDocument();
    });
  });

  it('should display cancel and confirm buttons in confirmation', async () => {
    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButton);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Cancel/i })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /Confirm Stop/i })).toBeInTheDocument();
    });
  });

  it('should return to form when cancel is clicked', async () => {
    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButton);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Cancel/i })).toBeInTheDocument();
    });

    const cancelButton = screen.getByRole('button', { name: /Cancel/i });
    fireEvent.click(cancelButton);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Stop Session/i })).toBeInTheDocument();
    });
  });

  it('should call stopSession API when confirm is clicked', async () => {
    apiClient.stopSession.mockResolvedValueOnce(undefined);

    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButton);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Confirm Stop/i })).toBeInTheDocument();
    });

    const confirmButton = screen.getByRole('button', { name: /Confirm Stop/i });
    fireEvent.click(confirmButton);

    await waitFor(() => {
      expect(apiClient.stopSession).toHaveBeenCalledWith('+79991234567');
    });
  });

  it('should call onSuccess callback after successful stop', async () => {
    apiClient.stopSession.mockResolvedValueOnce(undefined);

    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButton);

    await waitFor(() => {
      const confirmButton = screen.getByRole('button', { name: /Confirm Stop/i });
      fireEvent.click(confirmButton);
    });

    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalled();
    });
  });

  it('should show error message on stop failure', async () => {
    const errorMessage = 'Session not found';
    apiClient.stopSession.mockRejectedValueOnce(new Error(errorMessage));

    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButton);

    await waitFor(() => {
      const confirmButton = screen.getByRole('button', { name: /Confirm Stop/i });
      fireEvent.click(confirmButton);
    });

    await waitFor(() => {
      expect(screen.getByText(new RegExp(errorMessage, 'i'))).toBeInTheDocument();
    });
  });

  it('should reset form after successful stop', async () => {
    apiClient.stopSession.mockResolvedValueOnce(undefined);

    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButton);

    await waitFor(() => {
      const confirmButton = screen.getByRole('button', { name: /Confirm Stop/i });
      fireEvent.click(confirmButton);
    });

    await waitFor(() => {
      // After successful stop, form should be displayed again
      expect(screen.getByRole('button', { name: /Stop Session/i })).toBeInTheDocument();
    });
  });

  it('should display loading state during stop', async () => {
    apiClient.stopSession.mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve(), 100))
    );

    render(<StopSession onSuccess={mockOnSuccess} />);

    const phoneInput = screen.getByLabelText(/Phone Number/i);
    fireEvent.change(phoneInput, { target: { value: '+79991234567' } });

    const stopButton = screen.getByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButton);

    await waitFor(() => {
      const confirmButton = screen.getByRole('button', { name: /Confirm Stop/i });
      fireEvent.click(confirmButton);
    });

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Stopping/i })).toBeInTheDocument();
    });
  });
});
