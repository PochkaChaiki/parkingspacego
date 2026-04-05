import '@/styles/globals.css';

export const metadata = {
  title: 'Parking Management System',
  description: 'Manage your parking sessions',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
