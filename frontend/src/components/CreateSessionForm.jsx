'use client';

import React, { useState } from 'react';
import { createSession } from '@/lib/apiClient';
import styles from './CreateSessionForm.module.css';

export default function CreateSessionForm({ onSuccess }) {
  const [formData, setFormData] = useState({
    client_name: '',
    phone_number: '',
    license_plate: '',
    spot_number: '',
    duration: '',
  });

  const [error, setError] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    setError(null);
  };

  const isFormValid = () => {
    return (
      formData.client_name.trim() !== '' &&
      formData.phone_number.trim() !== '' &&
      formData.license_plate.trim() !== '' &&
      formData.spot_number.trim() !== ''
    );
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      // Convert spot_number to number
      const submitData = {
        ...formData,
        spot_number: parseInt(formData.spot_number, 10),
      };

      // Remove empty duration field
      if (!submitData.duration) {
        delete submitData.duration;
      }

      const response = await createSession(submitData);

      // Check response status and handle accordingly
      if (response.status === 'failure') {
        setError('Session with this phone number is already active');
        return;
      }

      if (response.status === 'occupied') {
        setError('Chosen spot is already occupied');
        return;
      }

      // Reset form on success
      setFormData({
        client_name: '',
        phone_number: '',
        license_plate: '',
        spot_number: '',
        duration: '',
      });

      // Call success callback
      if (onSuccess) {
        onSuccess();
      }
    } catch (err) {
      setError(err.message || 'Failed to create session');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <h2>Create Parking Session</h2>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.formGroup}>
        <label htmlFor="client_name">Client Name</label>
        <input
          id="client_name"
          type="text"
          name="client_name"
          value={formData.client_name}
          onChange={handleChange}
          placeholder="Enter client name"
          required
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="phone_number">Phone Number</label>
        <input
          id="phone_number"
          type="tel"
          name="phone_number"
          value={formData.phone_number}
          onChange={handleChange}
          placeholder="7 999 123 45 67"
          required
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="license_plate">License Plate</label>
        <input
          id="license_plate"
          type="text"
          name="license_plate"
          value={formData.license_plate}
          onChange={handleChange}
          placeholder="A123BC140"
          required
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="spot_number">Spot Number</label>
        <input
          id="spot_number"
          type="number"
          name="spot_number"
          value={formData.spot_number}
          onChange={handleChange}
          placeholder="42"
          min="1"
          required
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="duration">Duration (Optional)</label>
        <input
          id="duration"
          type="text"
          name="duration"
          value={formData.duration}
          onChange={handleChange}
          placeholder="1h"
        />
      </div>

      <button
        type="submit"
        disabled={!isFormValid() || isLoading}
        className={styles.submitButton}
      >
        {isLoading ? 'Creating...' : 'Create Session'}
      </button>
    </form>
  );
}
