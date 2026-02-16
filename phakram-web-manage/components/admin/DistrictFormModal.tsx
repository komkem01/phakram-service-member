'use client';

import { useState, useEffect } from 'react';
import { District } from '@/lib/api/districts';
import { Province, provincesApi } from '@/lib/api/provinces';
import SearchableSelect from '@/components/ui/SearchableSelect';

interface DistrictFormModalProps {
  isOpen: boolean;
  mode: 'create' | 'edit';
  district?: District | null;
  onClose: () => void;
  onSubmit: (data: { province_id: string; name: string; is_active: boolean }) => void;
}

export default function DistrictFormModal({
  isOpen,
  mode,
  district,
  onClose,
  onSubmit,
}: DistrictFormModalProps) {
  const [provinceId, setProvinceId] = useState('');
  const [name, setName] = useState('');
  const [isActive, setIsActive] = useState(true);
  const [provinces, setProvinces] = useState<Province[]>([]);
  const [loadingProvinces, setLoadingProvinces] = useState(false);

  useEffect(() => {
    if (isOpen) {
      loadProvinces();
      if (mode === 'edit' && district) {
        setProvinceId(district.province_id);
        setName(district.name);
        setIsActive(district.is_active);
      } else {
        setProvinceId('');
        setName('');
        setIsActive(true);
      }
    }
  }, [isOpen, mode, district]);

  const loadProvinces = async () => {
    try {
      setLoadingProvinces(true);
      const response = await provincesApi.list({ page: 1, size: 1000 });
      setProvinces(response.data);
    } catch (error) {
      console.error('Failed to load provinces:', error);
    } finally {
      setLoadingProvinces(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ province_id: provinceId, name, is_active: isActive });
  };

  if (!isOpen) return null;

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
              {mode === 'create' ? 'เพิ่มอำเภอใหม่' : 'แก้ไขอำเภอ'}
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
            {/* Province */}
            <div>
              <SearchableSelect
                label="จังหวัด"
                options={provinces.map(p => ({ id: p.id, name: p.name }))}
                value={provinceId}
                onChange={setProvinceId}
                placeholder="-- เลือกจังหวัด --"
                required
                disabled={loadingProvinces}
              />
            </div>

            {/* Name */}
            <div>
              <label
                htmlFor="district-name"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                ชื่ออำเภอ <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="district-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="กรุณากรอกชื่ออำเภอ"
                required
              />
            </div>

            {/* Is Active */}
            <div className="flex items-center">
              <input
                type="checkbox"
                id="district-active"
                checked={isActive}
                onChange={(e) => setIsActive(e.target.checked)}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <label
                htmlFor="district-active"
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
              {mode === 'create' ? 'เพิ่มอำเภอ' : 'บันทึก'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
