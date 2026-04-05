'use client';

import React, { useState } from 'react';
import { stopSession } from '@/lib/apiClient';
import styles from './StopSession.module.css';

export default function StopSession({ onSuccess }) {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [error, setError] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);

  const handleStop = async () => {
    setError(null);
    setIsLoading(true);

    try {
      await stopSession(phoneNumber.replaceAll(' ', ''));

      // Reset form
      setPhoneNumber('');
      setShowConfirm(false);

      if (onSuccess) {
        onSuccess();
      }
    } catch (err) {
      setError(err.message || 'Failed to stop session');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    setShowConfirm(true);
  };

  return (
    <div className={styles.container}>
      <h2>Stop Parking Session</h2>

      {error && <div className={styles.error}>{error}</div>}

      {!showConfirm ? (
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGroup}>
            <label htmlFor="phone_number">Phone Number</label>
            <input
              id="phone_number"
              type="tel"
              value={phoneNumber}
              onChange={(e) => setPhoneNumber(e.target.value)}
              placeholder="+7 999 123 45 67"
              required
            />
          </div>

          <button
            type="submit"
            disabled={phoneNumber.trim() === '' || isLoading}
            className={styles.stopButton}
          >
            {isLoading ? 'Processing...' : 'Stop Session'}
          </button>
        </form>
      ) : (
        <div className={styles.confirmation}>
          <p>Are you sure you want to stop the session for {phoneNumber}?</p>
          <div className={styles.buttonGroup}>
            <button
              onClick={() => setShowConfirm(false)}
              disabled={isLoading}
              className={styles.cancelButton}
            >
              Cancel
            </button>
            <button
              onClick={handleStop}
              disabled={isLoading}
              className={styles.confirmButton}
            >
              {isLoading ? 'Stopping...' : 'Confirm Stop'}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
