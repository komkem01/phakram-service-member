import type { Metadata } from "next";
import "./globals.css";
import { ModalProvider } from "@/contexts/ModalContext";

export const metadata: Metadata = {
  title: "Phakram Web Manage",
  description: "Management frontend for Phakram",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="th">
      <body className="antialiased">
        <ModalProvider>
          {children}
        </ModalProvider>
      </body>
    </html>
  );
}
