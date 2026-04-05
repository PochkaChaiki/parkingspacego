/**
 * Component Tests - Session Creation Form
 * Following Red/Green TDD methodology
 */
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import CreateSessionForm from '../CreateSessionForm';
import * as apiClient from '@/lib/apiClient';

jest.mock('@/lib/apiClient');

describe('CreateSessionForm Component', () => {
  const mockOnSuccess = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render form with all required fields', () => {
    render(<CreateSessionForm onSuccess={mockOnSuccess} />);

    expect(screen.getByLabelText(/Client Name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Phone Number/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/License Plate/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Spot Number/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Create Session/i })).toBeInTheDocument();
  });

  it('should have optional duration field', () => {
    render(<CreateSessionForm onSuccess={mockOnSuccess} />);
    expect(screen.getByLabelText(/Duration/i)).toBeInTheDocument();
  });

  it('should disable submit button initially', () => {
    render(<CreateSessionForm onSuccess={mockOnSuccess} />);
    const submitButton = screen.getByRole('button', { name: /Create Session/i });
    expect(submitButton).toBeDisabled();
  });

  it('should enable submit button when all required fields are filled', async () => {
    render(<CreateSessionForm onSuccess={mockOnSuccess} />);

    fireEvent.change(screen.getByLabelText(/Client Name/i), {
      target: { value: 'Иван' },
    });
    fireEvent.change(screen.getByLabelText(/Phone Number/i), {
      target: { value: '+79991234567' },
    });
    fireEvent.change(screen.getByLabelText(/License Plate/i), {
      target: { value: 'A123BC140' },
    });
    fireEvent.change(screen.getByLabelText(/Spot Number/i), {
      target: { value: '42' },
    });

    await waitFor(() => {
      const submitButton = screen.getByRole('button', { name: /Create Session/i });
      expect(submitButton).not.toBeDisabled();
    });
  });

  it('should call createSession API when form is submitted', async () => {
    apiClient.createSession.mockResolvedValueOnce({ status: 'success' });

    render(<CreateSessionForm onSuccess={mockOnSuccess} />);

    fireEvent.change(screen.getByLabelText(/Client Name/i), {
      target: { value: 'Иван' },
    });
    fireEvent.change(screen.getByLabelText(/Phone Number/i), {
      target: { value: '+79991234567' },
    });
    fireEvent.change(screen.getByLabelText(/License Plate/i), {
      target: { value: 'A123BC140' },
    });
    fireEvent.change(screen.getByLabelText(/Spot Number/i), {
      target: { value: '42' },
    });

    const submitButton = screen.getByRole('button', { name: /Create Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(apiClient.createSession).toHaveBeenCalled();
    });
  });

  it('should call onSuccess callback on successful submission', async () => {
    apiClient.createSession.mockResolvedValueOnce({ status: 'success' });

    render(<CreateSessionForm onSuccess={mockOnSuccess} />);

    fireEvent.change(screen.getByLabelText(/Client Name/i), {
      target: { value: 'Иван' },
    });
    fireEvent.change(screen.getByLabelText(/Phone Number/i), {
      target: { value: '+79991234567' },
    });
    fireEvent.change(screen.getByLabelText(/License Plate/i), {
      target: { value: 'A123BC140' },
    });
    fireEvent.change(screen.getByLabelText(/Spot Number/i), {
      target: { value: '42' },
    });

    const submitButton = screen.getByRole('button', { name: /Create Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalled();
    });
  });

  it('should show error message on submission failure', async () => {
    const errorMessage = 'Network error';
    apiClient.createSession.mockRejectedValueOnce(new Error(errorMessage));

    render(<CreateSessionForm onSuccess={mockOnSuccess} />);

    fireEvent.change(screen.getByLabelText(/Client Name/i), {
      target: { value: 'Иван' },
    });
    fireEvent.change(screen.getByLabelText(/Phone Number/i), {
      target: { value: '+79991234567' },
    });
    fireEvent.change(screen.getByLabelText(/License Plate/i), {
      target: { value: 'A123BC140' },
    });
    fireEvent.change(screen.getByLabelText(/Spot Number/i), {
      target: { value: '42' },
    });

    const submitButton = screen.getByRole('button', { name: /Create Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(new RegExp(errorMessage, 'i'))).toBeInTheDocument();
    });
  });

  it('should reset form after successful submission', async () => {
    apiClient.createSession.mockResolvedValueOnce({ status: 'success' });

    render(<CreateSessionForm onSuccess={mockOnSuccess} />);

    fireEvent.change(screen.getByLabelText(/Client Name/i), {
      target: { value: 'Иван' },
    });
    fireEvent.change(screen.getByLabelText(/Phone Number/i), {
      target: { value: '+79991234567' },
    });
    fireEvent.change(screen.getByLabelText(/License Plate/i), {
      target: { value: 'A123BC140' },
    });
    fireEvent.change(screen.getByLabelText(/Spot Number/i), {
      target: { value: '42' },
    });

    const submitButton = screen.getByRole('button', { name: /Create Session/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect((screen.getByLabelText(/Client Name/i) ).value).toBe('');
      expect((screen.getByLabelText(/Phone Number/i)).value).toBe('');
    });
  });
});
