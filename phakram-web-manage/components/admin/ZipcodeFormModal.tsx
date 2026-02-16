'use client';

import { useState, useEffect } from 'react';
import { Zipcode } from '@/lib/api/zipcodes';
import { SubDistrict, subDistrictsApi } from '@/lib/api/sub-districts';
import SearchableSelect from '@/components/ui/SearchableSelect';

interface ZipcodeFormModalProps {
  isOpen: boolean;
  mode: 'create' | 'edit';
  zipcode?: Zipcode | null;
  onClose: () => void;
  onSubmit: (data: { sub_districts_id: string; name: string; is_active: boolean }) => void;
}

export default function ZipcodeFormModal({
  isOpen,
  mode,
  zipcode,
  onClose,
  onSubmit,
}: ZipcodeFormModalProps) {
  const [subDistrictId, setSubDistrictId] = useState('');
  const [name, setName] = useState('');
  const [isActive, setIsActive] = useState(true);
  const [subDistricts, setSubDistricts] = useState<SubDistrict[]>([]);
  const [loadingSubDistricts, setLoadingSubDistricts] = useState(false);

  useEffect(() => {
    if (isOpen) {
      loadSubDistricts();
      if (mode === 'edit' && zipcode) {
        setSubDistrictId(zipcode.sub_districts_id ?? zipcode.sub_district_id ?? '');
        setName(zipcode.name);
        setIsActive(zipcode.is_active);
      } else {
        setSubDistrictId('');
        setName('');
        setIsActive(true);
      }
    }
  }, [isOpen, mode, zipcode]);

  const loadSubDistricts = async () => {
    try {
      setLoadingSubDistricts(true);
      const response = await subDistrictsApi.list({ page: 1, size: 1000 });
      setSubDistricts(response.data);
    } catch (error) {
      console.error('Failed to load sub-districts:', error);
    } finally {
      setLoadingSubDistricts(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ sub_districts_id: subDistrictId, name, is_active: isActive });
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
              {mode === 'create' ? 'เพิ่มรหัสไปรษณีย์ใหม่' : 'แก้ไขรหัสไปรษณีย์'}
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
            {/* Sub District */}
            <div>
              <SearchableSelect
                label="ตำบล"
                options={subDistricts.map(sd => ({ id: sd.id, name: sd.name }))}
                value={subDistrictId}
                onChange={setSubDistrictId}
                placeholder="-- เลือกตำบล --"
                required
                disabled={loadingSubDistricts}
              />
            </div>

            {/* Name */}
            <div>
              <label
                htmlFor="zipcode-name"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                รหัสไปรษณีย์ <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="zipcode-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="กรุณากรอกรหัสไปรษณีย์ (5 หลัก)"
                maxLength={5}
                pattern="[0-9]{5}"
                required
              />
            </div>

            {/* Is Active */}
            <div className="flex items-center">
              <input
                type="checkbox"
                id="zipcode-active"
                checked={isActive}
                onChange={(e) => setIsActive(e.target.checked)}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <label
                htmlFor="zipcode-active"
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
              {mode === 'create' ? 'เพิ่มรหัสไปรษณีย์' : 'บันทึก'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
