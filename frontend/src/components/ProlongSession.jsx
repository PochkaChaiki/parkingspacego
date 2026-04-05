'use client';

import React, { useState } from 'react';
import { prolongSession } from '@/lib/apiClient';
import styles from './ProlongSession.module.css';

export default function ProlongSession({ onSuccess }) {
  const [formData, setFormData] = useState({
    phone_number: '',
    duration: '1h',
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

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      await prolongSession(formData.phone_number.replaceAll(' ', ''), formData.duration);

      // Reset form
      setFormData({
        phone_number: '',
        duration: '1h',
      });

      if (onSuccess) {
        onSuccess();
      }
    } catch (err) {
      setError(err.message || 'Failed to prolong session');
    } finally {
      setIsLoading(false);
    }
  };

  const isFormValid = () => {
    return formData.phone_number.trim() !== '' && formData.duration.trim() !== '';
  };

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <h2>Prolong Parking Session</h2>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.formGroup}>
        <label htmlFor="phone_number">Phone Number</label>
        <input
          id="phone_number"
          type="tel"
          name="phone_number"
          value={formData.phone_number}
          onChange={handleChange}
          placeholder="+79991234567"
          required
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="duration">Duration to Add</label>
        <input
          id="duration"
          type="text"
          name="duration"
          value={formData.duration}
          onChange={handleChange}
          placeholder="1h"
          required
        />
      </div>

      <button
        type="submit"
        disabled={!isFormValid() || isLoading}
        className={styles.submitButton}
      >
        {isLoading ? 'Prolonging...' : 'Prolong Session'}
      </button>
    </form>
  );
}
