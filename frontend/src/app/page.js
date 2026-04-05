'use client';

import React, { useState } from 'react';
import CreateSessionForm from '@/components/CreateSessionForm';
import GetSession from '@/components/GetSession';
import ProlongSession from '@/components/ProlongSession';
import StopSession from '@/components/StopSession';
import styles from './page.module.css';

export default function Home() {
  const [activeTab, setActiveTab] = useState('create');
  const [successMessage, setSuccessMessage] = useState('');

  const handleSuccess = () => {
    setSuccessMessage('Operation completed successfully!');
    setTimeout(() => setSuccessMessage(''), 3000);
  };


  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1>Parking Management System</h1>
        <p>Manage your parking sessions easily</p>
      </header>

      {successMessage && <div className={styles.successMessage}>{successMessage}</div>}

      <nav className={styles.tabs}>
        <button
          className={`${styles.tabButton} ${activeTab === 'create' ? styles.active : ''}`}
          onClick={() => setActiveTab('create')}
        >
          Create Session
        </button>
        <button
          className={`${styles.tabButton} ${activeTab === 'view' ? styles.active : ''}`}
          onClick={() => setActiveTab('view')}
        >
          View Session
        </button>
        <button
          className={`${styles.tabButton} ${activeTab === 'prolong' ? styles.active : ''}`}
          onClick={() => setActiveTab('prolong')}
        >
          Prolong Session
        </button>
        <button
          className={`${styles.tabButton} ${activeTab === 'stop' ? styles.active : ''}`}
          onClick={() => setActiveTab('stop')}
        >
          Stop Session
        </button>
      </nav>

      <div className={styles.content}>
        {activeTab === 'create' && <CreateSessionForm onSuccess={handleSuccess} />}
        {activeTab === 'view' && <GetSession />}
        {activeTab === 'prolong' && <ProlongSession onSuccess={handleSuccess} />}
        {activeTab === 'stop' && <StopSession onSuccess={handleSuccess} />}
      </div>

      <footer className={styles.footer}>
        <p>
          🔗 Backend API: <code>{process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}</code>
        </p>
      </footer>
    </div>
  );
}
