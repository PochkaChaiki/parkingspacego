import React from 'react';
import './globals.css';
import PropTypes from 'prop-types';

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

RootLayout.propTypes = {
  children: PropTypes.node.isRequired,
};
