import React from 'react';
import './globals.css';

export const metadata = {
  title: 'Parking Management System',
  description: 'Frontend for parking management system',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
