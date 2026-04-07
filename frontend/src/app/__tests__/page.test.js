/**
 * Integration Tests - Main Page
 * Testing the overall parking management system interface
 */
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import Home from '../page';
import * as apiClient from '@/lib/apiClient';

jest.mock('@/lib/apiClient');

describe('Parking Management System - Main Page', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render main header with title', () => {
    render(<Home />);
    expect(screen.getByRole('heading', { name: /Parking Management System/i, level: 1 })).toBeInTheDocument();
    expect(screen.getByText(/Manage your parking sessions/i)).toBeInTheDocument();
  });

  it('should render all four navigation tabs', () => {
    render(<Home />);
    // Check that tabs exist (getByRole will fail if multiple exist, so use queryAll)
    expect(screen.queryAllByRole('button', { name: /Create Session/i }).length).toBeGreaterThan(0);
    expect(screen.queryAllByRole('button', { name: /View Session/i }).length).toBeGreaterThan(0);
    expect(screen.queryAllByRole('button', { name: /Prolong Session/i }).length).toBeGreaterThan(0);
    expect(screen.queryAllByRole('button', { name: /Stop Session/i }).length).toBeGreaterThan(0);
  });

  it('should display Create Session form by default', () => {
    render(<Home />);
    expect(screen.getByText(/Create Parking Session/i)).toBeInTheDocument();
  });

  it('should switch to View Session tab', async () => {
    render(<Home />);
    const viewButtons = screen.getAllByRole('button', { name: /View Session/i });
    fireEvent.click(viewButtons[0]); // Click the tab button

    await waitFor(() => {
      expect(screen.getByText(/View Parking Session/i)).toBeInTheDocument();
    });
  });

  it('should switch to Prolong Session tab', async () => {
    render(<Home />);
    const prolongButtons = screen.getAllByRole('button', { name: /Prolong Session/i });
    fireEvent.click(prolongButtons[0]); // Click the tab button

    await waitFor(() => {
      expect(screen.getByText(/Prolong Parking Session/i)).toBeInTheDocument();
    });
  });

  it('should switch to Stop Session tab', async () => {
    render(<Home />);
    const stopButtons = screen.getAllByRole('button', { name: /Stop Session/i });
    fireEvent.click(stopButtons[0]); // Click the tab button

    await waitFor(() => {
      expect(screen.getByText(/Stop Parking Session/i)).toBeInTheDocument();
    });
  });

  it('should show active tab styling', async () => {
    const { container } = render(<Home />);
    const tabButtons = screen.getAllByRole('button', { name: /Create Session/i });
    const createTabButton = tabButtons[0]; // First one is the tab button

    expect(createTabButton).toHaveClass('active');

    const viewTabButtons = screen.getAllByRole('button', { name: /View Session/i });
    const viewTabButton = viewTabButtons[0];
    fireEvent.click(viewTabButton);

    await waitFor(() => {
      expect(createTabButton).not.toHaveClass('active');
      expect(viewTabButton).toHaveClass('active');
    });
  });

  it('should call API when form is submitted', async () => {
    apiClient.createSession.mockResolvedValueOnce({ status: 'success' });

    render(<Home />);

    // Fill in the form
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

    const buttons = screen.getAllByRole('button', { name: /Create Session/i });
    fireEvent.click(buttons[buttons.length - 1]); // Click the form button, not the tab

    await waitFor(() => {
      expect(apiClient.createSession).toHaveBeenCalled();
    });
  });

  it('should display success message on successful operation', async () => {
    apiClient.prolongSession.mockResolvedValueOnce({ status: 'success' });

    render(<Home />);

    // Switch to Prolong tab
    const tabButtons = screen.getAllByRole('button', { name: /Prolong Session/i });
    fireEvent.click(tabButtons[0]); // Click the tab button

    await waitFor(() => {
      expect(screen.getByText(/Prolong Parking Session/i)).toBeInTheDocument();
    });

    // Fill in the form
    fireEvent.change(screen.getByLabelText(/Phone Number/i), {
      target: { value: '+79991234567' },
    });

    const formButtons = screen.getAllByRole('button', { name: /Prolong Session/i });
    fireEvent.click(formButtons[formButtons.length - 1]); // Click the form button

    await waitFor(() => {
      expect(screen.getByText(/Session prolonged successfully/i)).toBeInTheDocument();
    });
  });

  it('should render footer with copyright', () => {
    render(<Home />);
    expect(screen.getByText(/© 2026 Parking Management System/i)).toBeInTheDocument();
  });

  it('should maintain tab state when switching between tabs', async () => {
    render(<Home />);

    // Switch to View tab
    const viewButton = screen.getByRole('button', { name: /View Session/i });
    fireEvent.click(viewButton);

    await waitFor(() => {
      expect(screen.getByText(/View Parking Session/i)).toBeInTheDocument();
    });

    // Switch to Prolong tab
    const prolongButton = screen.getByRole('button', { name: /Prolong Session/i });
    fireEvent.click(prolongButton);

    await waitFor(() => {
      expect(screen.getByText(/Prolong Parking Session/i)).toBeInTheDocument();
    });

    // Switch back to View tab
    fireEvent.click(viewButton);

    await waitFor(() => {
      expect(screen.getByText(/View Parking Session/i)).toBeInTheDocument();
    });
  });
});
