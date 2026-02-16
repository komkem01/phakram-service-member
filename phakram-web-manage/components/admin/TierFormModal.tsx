'use client';

import { useEffect, useState } from 'react';

interface Tier {
  id: string;
  name_th: string;
  name_en: string;
  min_spending: number;
  discount_rate: number;
  is_active: boolean;
}

interface TierFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: {
    name_th: string;
    name_en: string;
    min_spending: string;
    discount_rate: string;
    is_active: boolean;
  }) => void;
  tier?: Tier | null;
  mode: 'create' | 'edit';
}

export default function TierFormModal({
  isOpen,
  onClose,
  onSubmit,
  tier,
  mode,
}: TierFormModalProps) {
  const [formData, setFormData] = useState({
    name_th: '',
    name_en: '',
    min_spending: '',
    discount_rate: '',
    is_active: true,
  });

  useEffect(() => {
    if (isOpen) {
      if (mode === 'edit' && tier) {
        setFormData({
          name_th: tier.name_th,
          name_en: tier.name_en,
          min_spending: tier.min_spending.toString(),
          discount_rate: tier.discount_rate.toString(),
          is_active: tier.is_active,
        });
      } else {
        setFormData({
          name_th: '',
          name_en: '',
          min_spending: '',
          discount_rate: '',
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
  }, [isOpen, mode, tier]);

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
              {mode === 'create' ? 'เพิ่มระดับสมาชิกใหม่' : 'แก้ไขระดับสมาชิก'}
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
                ชื่อระดับสมาชิก (ไทย) <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="name_th"
                name="name_th"
                value={formData.name_th}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="เช่น ทอง, เงิน, ทองแดง"
              />
            </div>

            {/* Name English */}
            <div>
              <label
                htmlFor="name_en"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                ชื่อระดับสมาชิก (อังกฤษ) <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="name_en"
                name="name_en"
                value={formData.name_en}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="e.g. Gold, Silver, Bronze"
              />
            </div>

            {/* Min Spending */}
            <div>
              <label
                htmlFor="min_spending"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                ยอดใช้จ่ายขั้นต่ำ (บาท) <span className="text-red-500">*</span>
              </label>
              <input
                type="number"
                id="min_spending"
                name="min_spending"
                value={formData.min_spending}
                onChange={handleChange}
                required
                step="0.01"
                min="0"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="0.00"
              />
            </div>

            {/* Discount Rate */}
            <div>
              <label
                htmlFor="discount_rate"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                อัตราส่วนลด (%) <span className="text-red-500">*</span>
              </label>
              <input
                type="number"
                id="discount_rate"
                name="discount_rate"
                value={formData.discount_rate}
                onChange={handleChange}
                required
                step="0.01"
                min="0"
                max="100"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="0.00"
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
              {mode === 'create' ? 'เพิ่มระดับสมาชิก' : 'บันทึก'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
