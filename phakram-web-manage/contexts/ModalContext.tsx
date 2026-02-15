'use client';

import { createContext, useContext, useState, ReactNode } from 'react';
import Modal, { ModalType } from '@/components/Modal';

interface ModalOptions {
  title: string;
  message: string;
  type?: ModalType;
  confirmText?: string;
  cancelText?: string;
  onConfirm?: () => void;
}

interface ModalContextType {
  showModal: (options: ModalOptions) => void;
  showSuccess: (message: string, title?: string) => void;
  showError: (message: string, title?: string) => void;
  showConfirm: (message: string, onConfirm: () => void, title?: string) => void;
  showWarning: (message: string, onConfirm: () => void, title?: string) => void;
  closeModal: () => void;
}

const ModalContext = createContext<ModalContextType | undefined>(undefined);

export function ModalProvider({ children }: { children: ReactNode }) {
  const [isOpen, setIsOpen] = useState(false);
  const [modalOptions, setModalOptions] = useState<ModalOptions>({
    title: '',
    message: '',
    type: 'info',
  });

  const showModal = (options: ModalOptions) => {
    setModalOptions(options);
    setIsOpen(true);
  };

  const closeModal = () => {
    setIsOpen(false);
  };

  const showSuccess = (message: string, title = 'สำเร็จ') => {
    showModal({ title, message, type: 'success' });
  };

  const showError = (message: string, title = 'เกิดข้อผิดพลาด') => {
    showModal({ title, message, type: 'error' });
  };

  const showConfirm = (
    message: string,
    onConfirm: () => void,
    title = 'ยืนยันการดำเนินการ'
  ) => {
    showModal({
      title,
      message,
      type: 'confirm',
      onConfirm,
    });
  };

  const showWarning = (
    message: string,
    onConfirm: () => void,
    title = 'คำเตือน'
  ) => {
    showModal({
      title,
      message,
      type: 'warning',
      onConfirm,
      confirmText: 'ลบ',
      cancelText: 'ยกเลิก',
    });
  };

  return (
    <ModalContext.Provider
      value={{ showModal, showSuccess, showError, showConfirm, showWarning, closeModal }}
    >
      {children}
      <Modal
        isOpen={isOpen}
        onClose={closeModal}
        title={modalOptions.title}
        message={modalOptions.message}
        type={modalOptions.type}
        onConfirm={modalOptions.onConfirm}
        confirmText={modalOptions.confirmText}
        cancelText={modalOptions.cancelText}
      />
    </ModalContext.Provider>
  );
}

export function useModal() {
  const context = useContext(ModalContext);
  if (context === undefined) {
    throw new Error('useModal must be used within a ModalProvider');
  }
  return context;
}
