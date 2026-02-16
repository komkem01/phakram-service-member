'use client';

import { useEffect, useState } from 'react';

interface Bank {
  id: string;
  name_th: string;
  name_abb_th: string;
  name_en: string;
  name_abb_en: string;
  is_active: boolean;
}

interface BankFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: {
    name_th: string;
    name_abb_th: string;
    name_en: string;
    name_abb_en: string;
    is_active: boolean;
  }) => void;
  bank?: Bank | null;
  mode: 'create' | 'edit';
}

export default function BankFormModal({
  isOpen,
  onClose,
  onSubmit,
  bank,
  mode,
}: BankFormModalProps) {
  const [formData, setFormData] = useState({
    name_th: '',
    name_abb_th: '',
    name_en: '',
    name_abb_en: '',
    is_active: true,
  });

  useEffect(() => {
    if (isOpen) {
      if (mode === 'edit' && bank) {
        setFormData({
          name_th: bank.name_th,
          name_abb_th: bank.name_abb_th,
          name_en: bank.name_en,
          name_abb_en: bank.name_abb_en,
          is_active: bank.is_active,
        });
      } else {
        setFormData({
          name_th: '',
          name_abb_th: '',
          name_en: '',
          name_abb_en: '',
          is_active: true,
        });
      }
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isOpen, mode, bank]);

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
    onClose();
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative bg-white rounded-xl shadow-2xl w-full max-w-md mx-4 animate-modal-show">
        <form onSubmit={handleSubmit}>
          {/* Header */}
          <div className="flex items-center justify-between p-6 border-b border-gray-200">
            <h3 className="text-xl font-semibold text-gray-900">
              {mode === 'create' ? 'เพิ่มธนาคารใหม่' : 'แก้ไขธนาคาร'}
            </h3>
            <button
              type="button"
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg
                className="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>

          {/* Body */}
          <div className="p-6 space-y-4">
            {/* Name Thai */}
            <div>
              <label
                htmlFor="name_th"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                ชื่อธนาคาร (ไทย) <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="name_th"
                name="name_th"
                value={formData.name_th}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="เช่น ธนาคารกรุงเทพ"
              />
            </div>

            {/* Name Abbreviation Thai */}
            <div>
              <label
                htmlFor="name_abb_th"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                ชื่อย่อ (ไทย) <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="name_abb_th"
                name="name_abb_th"
                value={formData.name_abb_th}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="เช่น กรุงเทพ"
              />
            </div>

            {/* Name English */}
            <div>
              <label
                htmlFor="name_en"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                ชื่อธนาคาร (อังกฤษ) <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="name_en"
                name="name_en"
                value={formData.name_en}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="e.g. Bangkok Bank"
              />
            </div>

            {/* Name Abbreviation English */}
            <div>
              <label
                htmlFor="name_abb_en"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                ชื่อย่อ (อังกฤษ) <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="name_abb_en"
                name="name_abb_en"
                value={formData.name_abb_en}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="e.g. BBL"
              />
            </div>

            {/* Is Active */}
            <div className="flex items-center">
              <input
                type="checkbox"
                id="is_active"
                name="is_active"
                checked={formData.is_active}
                onChange={handleChange}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <label
                htmlFor="is_active"
                className="ml-2 text-sm font-medium text-gray-700"
              >
                เปิดใช้งาน
              </label>
            </div>
          </div>

          {/* Footer */}
          <div className="flex items-center justify-end gap-3 p-6 border-t border-gray-200 bg-gray-50 rounded-b-xl">
            <button
              type="button"
              onClick={onClose}
              className="px-6 py-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors font-medium"
            >
              ยกเลิก
            </button>
            <button
              type="submit"
              className="px-6 py-2 text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors font-medium"
            >
              {mode === 'create' ? 'เพิ่มธนาคาร' : 'บันทึก'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
