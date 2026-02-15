'use client';

import { useEffect } from 'react';

export type ModalType = 'success' | 'error' | 'confirm' | 'warning' | 'info';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm?: () => void;
  title: string;
  message: string;
  type?: ModalType;
  confirmText?: string;
  cancelText?: string;
}

const modalStyles: Record<ModalType, { icon: string; iconBg: string; iconColor: string }> = {
  success: {
    icon: '✓',
    iconBg: 'bg-emerald-100',
    iconColor: 'text-emerald-600',
  },
  error: {
    icon: '✕',
    iconBg: 'bg-red-100',
    iconColor: 'text-red-600',
  },
  warning: {
    icon: '!',
    iconBg: 'bg-amber-100',
    iconColor: 'text-amber-600',
  },
  confirm: {
    icon: '?',
    iconBg: 'bg-blue-100',
    iconColor: 'text-blue-600',
  },
  info: {
    icon: 'i',
    iconBg: 'bg-slate-100',
    iconColor: 'text-slate-600',
  },
};

export default function Modal({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  type = 'info',
  confirmText = 'ยืนยัน',
  cancelText = 'ยกเลิก',
}: ModalProps) {
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isOpen]);

  if (!isOpen) return null;

  const style = modalStyles[type];
  const showActions = type === 'confirm' || type === 'warning';

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative bg-white rounded-xl shadow-2xl w-full max-w-md mx-4 animate-modal-show">
        <div className="p-6">
          {/* Icon */}
          <div className="flex justify-center mb-4">
            <div
              className={`w-16 h-16 rounded-full ${style.iconBg} flex items-center justify-center`}
            >
              <span className={`text-3xl font-bold ${style.iconColor}`}>
                {style.icon}
              </span>
            </div>
          </div>

          {/* Title */}
          <h3 className="text-xl font-semibold text-slate-900 text-center mb-2">
            {title}
          </h3>

          {/* Message */}
          <p className="text-slate-600 text-center text-sm leading-relaxed">
            {message}
          </p>
        </div>

        {/* Actions */}
        <div className="border-t border-slate-100 p-4 flex gap-3">
          {showActions && onConfirm ? (
            <>
              <button
                onClick={onClose}
                className="flex-1 px-4 py-2.5 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg transition-colors"
              >
                {cancelText}
              </button>
              <button
                onClick={() => {
                  onConfirm();
                  onClose();
                }}
                className={`flex-1 px-4 py-2.5 text-sm font-medium text-white rounded-lg transition-colors ${
                  type === 'warning'
                    ? 'bg-red-600 hover:bg-red-700'
                    : 'bg-blue-600 hover:bg-blue-700'
                }`}
              >
                {confirmText}
              </button>
            </>
          ) : (
            <button
              onClick={onClose}
              className="w-full px-4 py-2.5 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
            >
              ตกลง
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
