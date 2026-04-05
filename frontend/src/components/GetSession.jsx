'use client';

import React, { useState } from 'react';
import { getSession, formatDateTimeForDisplay } from '@/lib/apiClient';
import styles from './GetSession.module.css';

export default function GetSession() {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [session, setSession] = useState(null);
  const [error, setError] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSearch = async (e) => {
    e.preventDefault();
    setError(null);
    setSession(null);
    setIsLoading(true);

    try {
      const result = await getSession(phoneNumber);
      setSession(result);
    } catch (err) {
      setError(err.message || 'Failed to get session');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <h2>View Parking Session</h2>

      <form onSubmit={handleSearch} className={styles.form}>
        <div className={styles.formGroup}>
          <label htmlFor="phone_number">Phone Number</label>
          <input
            id="phone_number"
            type="tel"
            value={phoneNumber}
            onChange={(e) => setPhoneNumber(e.target.value)}
            placeholder="7 999 123 45 67"
            required
          />
        </div>
        <button type="submit" disabled={isLoading} className={styles.searchButton}>
          {isLoading ? 'Searching...' : 'Search'}
        </button>
      </form>

      {error && <div className={styles.error}>{error}</div>}

      {session && (
        <div className={styles.sessionInfo}>
          <h3>Session Information</h3>
          <div className={styles.infoGrid}>
            <div className={styles.infoItem}>
              <span className={styles.label}>Client Name:</span>
              <span className={styles.value}>{session.client_name}</span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.label}>Phone Number:</span>
              <span className={styles.value}>{session.phone_number}</span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.label}>License Plate:</span>
              <span className={styles.value}>{session.license_plate}</span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.label}>Spot Number:</span>
              <span className={styles.value}>{session.spot_number}</span>
            </div>
            {session.start_time && (
              <div className={styles.infoItem}>
                <span className={styles.label}>Started:</span>
                <span className={styles.value}>
                  {formatDateTimeForDisplay(session.start_time)}
                </span>
              </div>
            )}
            {session.end_time && (
              <div className={styles.infoItem}>
                <span className={styles.label}>Ends:</span>
                <span className={styles.value}>
                  {formatDateTimeForDisplay(session.end_time)}
                </span>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
