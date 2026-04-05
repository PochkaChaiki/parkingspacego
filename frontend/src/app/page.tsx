'use client';

import { useState } from 'react';
import CreateSessionForm from '@/components/CreateSessionForm';
import GetSession from '@/components/GetSession';
import ProlongSession from '@/components/ProlongSession';
import StopSession from '@/components/StopSession';
import styles from './page.module.css';

type Tab = 'create' | 'view' | 'prolong' | 'stop';

export default function Home() {
  const [activeTab, setActiveTab] = useState<Tab>('create');
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleSuccess = (message: string) => {
    setSuccessMessage(message);
    setTimeout(() => setSuccessMessage(null), 3000);
  };

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <header className={styles.header}>
          <h1>Parking Management System</h1>
          <p>Manage your parking sessions easily</p>
        </header>

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

        {successMessage && (
          <div className={styles.successMessage}>
            {successMessage}
          </div>
        )}

        <div className={styles.content}>
          {activeTab === 'create' && (
            <CreateSessionForm onSuccess={() => handleSuccess('Session created successfully!')} />
          )}
          {activeTab === 'view' && <GetSession />}
          {activeTab === 'prolong' && (
            <ProlongSession onSuccess={() => handleSuccess('Session prolonged successfully!')} />
          )}
          {activeTab === 'stop' && (
            <StopSession onSuccess={() => handleSuccess('Session stopped successfully!')} />
          )}
        </div>

        <footer className={styles.footer}>
          <p>
            Backend API:{' '}
            <code>{process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}</code>
          </p>
        </footer>
      </div>
    </main>
  );
}
