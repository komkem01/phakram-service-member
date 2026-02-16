'use client';

import { useEffect, useState } from 'react';
import { gendersApi, Gender } from '@/lib/api/genders';
import SearchableSelect from '@/components/ui/SearchableSelect';

interface Prefix {
  id: string;
  name_th: string;
  name_en: string;
  gender_id: string;
  is_active: boolean;
}

interface PrefixFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: Omit<Prefix, 'id'>) => void;
  prefix?: Prefix | null;
  mode: 'create' | 'edit';
}

export default function PrefixFormModal({
  isOpen,
  onClose,
  onSubmit,
  prefix,
  mode,
}: PrefixFormModalProps) {
  const [formData, setFormData] = useState({
    name_th: '',
    name_en: '',
    gender_id: '',
    is_active: true,
  });
  const [genders, setGenders] = useState<Gender[]>([]);
  const [isLoadingGenders, setIsLoadingGenders] = useState(false);

  useEffect(() => {
    if (isOpen) {
      // Load genders for dropdown
      const loadGenders = async () => {
        try {
          setIsLoadingGenders(true);
          const response = await gendersApi.list();
          setGenders(response.data.filter((g) => g.is_active));
        } catch (error) {
          console.error('Failed to load genders:', error);
        } finally {
          setIsLoadingGenders(false);
        }
      };

      loadGenders();

      // Set form data
      if (mode === 'edit' && prefix) {
        setFormData({
          name_th: prefix.name_th,
          name_en: prefix.name_en,
          gender_id: prefix.gender_id,
          is_active: prefix.is_active,
        });
      } else {
        setFormData({
          name_th: '',
          name_en: '',
          gender_id: '',
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
  }, [isOpen, mode, prefix]);

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
    onClose();
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    const checked = (e.target as HTMLInputElement).checked;
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
              {mode === 'create' ? 'เพิ่มคำนำหน้าใหม่' : 'แก้ไขคำนำหน้า'}
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
                คำนำหน้า (ไทย) <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="name_th"
                name="name_th"
                value={formData.name_th}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="เช่น นาย, นาง"
              />
            </div>

            {/* Name English */}
            <div>
              <label
                htmlFor="name_en"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                คำนำหน้า (อังกฤษ) <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="name_en"
                name="name_en"
                value={formData.name_en}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="e.g. Mr., Mrs."
              />
            </div>

            {/* Gender */}
            <div>
              <SearchableSelect
                label="เพศ"
                options={genders.map(g => ({ 
                  id: g.id, 
                  name: `${g.name_th} (${g.name_en})` 
                }))}
                value={formData.gender_id}
                onChange={(value) => setFormData(prev => ({ ...prev, gender_id: value }))}
                placeholder="-- เลือกเพศ --"
                required
                disabled={isLoadingGenders}
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
              {mode === 'create' ? 'เพิ่มคำนำหน้า' : 'บันทึก'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
