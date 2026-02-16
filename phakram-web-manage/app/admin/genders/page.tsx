"use client";

import { useState, useEffect, useCallback } from "react";
import Loading from "@/components/admin/Loading";
import GenderFormModal from "@/components/admin/GenderFormModal";
import { useModal } from "@/contexts/ModalContext";
import { gendersApi, Gender, CreateGenderDto, UpdateGenderDto } from "@/lib/api/genders";
import { formatDateOnly } from "@/lib/utils/date";

export default function GendersPage() {
  const [genders, setGenders] = useState<Gender[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [selectedGender, setSelectedGender] = useState<Gender | null>(null);
  const { showSuccess, showWarning, showError } = useModal();

  // Fetch genders from API
  const fetchGenders = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await gendersApi.list();
      setGenders(response.data);
    } catch (error) {
      console.error('Failed to fetch genders:', error);
      showError('ไม่สามารถโหลดข้อมูลเพศได้', 'เกิดข้อผิดพลาด');
    } finally {
      setIsLoading(false);
    }
  }, [showError]);

  useEffect(() => {
    fetchGenders();
  }, [fetchGenders]);

  const handleCreate = () => {
    setModalMode('create');
    setSelectedGender(null);
    setIsModalOpen(true);
  };

  const handleEdit = (gender: Gender) => {
    setModalMode('edit');
    setSelectedGender(gender);
    setIsModalOpen(true);
  };

  const handleDelete = (gender: Gender) => {
    showWarning(
      `คุณต้องการลบเพศ "${gender.name_th}" ใช่หรือไม่?`,
      async () => {
        try {
          await gendersApi.delete(gender.id);
          showSuccess(`ลบเพศ "${gender.name_th}" สำเร็จ`);
          // Refresh list after delete
          fetchGenders();
        } catch (error) {
          console.error('Failed to delete gender:', error);
          showError(`ไม่สามารถลบเพศ "${gender.name_th}" ได้`, 'เกิดข้อผิดพลาด');
        }
      },
      'ยืนยันการลบ'
    );
  };

  const handleSubmit = async (data: Omit<Gender, 'id' | 'created_at'>) => {
    try {
      if (modalMode === 'create') {
        const createData: CreateGenderDto = {
          name_th: data.name_th,
          name_en: data.name_en,
          is_active: data.is_active,
        };
        await gendersApi.create(createData);
        showSuccess(`เพิ่มเพศ "${data.name_th}" สำเร็จ`);
      } else if (modalMode === 'edit' && selectedGender) {
        const updateData: UpdateGenderDto = {
          name_th: data.name_th,
          name_en: data.name_en,
          is_active: data.is_active,
        };
        await gendersApi.update(selectedGender.id, updateData);
        showSuccess(`แก้ไขเพศ "${data.name_th}" สำเร็จ`);
      }
      // Refresh list after create/update
      fetchGenders();
    } catch (error) {
      console.error('Failed to save gender:', error);
      showError(
        `ไม่สามารถ${modalMode === 'create' ? 'เพิ่ม' : 'แก้ไข'}เพศได้`,
        'เกิดข้อผิดพลาด'
      );
    }
  };

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการเพศ</h1>
        <button 
          onClick={handleCreate}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          + เพิ่มเพศ
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50">
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ลำดับ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ชื่อเพศ (ไทย)</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ชื่อเพศ (อังกฤษ)</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">สถานะ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">วันที่สร้าง</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              {genders.map((gender, index) => (
                <tr key={gender.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-4 px-6">{index + 1}</td>
                  <td className="py-4 px-6 font-medium">{gender.name_th}</td>
                  <td className="py-4 px-6 text-gray-600">{gender.name_en}</td>
                  <td className="py-4 px-6">
                    <div className="flex justify-center">
                      <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                        gender.is_active 
                          ? 'bg-green-100 text-green-700' 
                          : 'bg-gray-100 text-gray-600'
                      }`}>
                        {gender.is_active ? 'เปิดใช้งาน' : 'ปิดใช้งาน'}
                      </span>
                    </div>
                  </td>
                  <td className="py-4 px-6 text-gray-600">{formatDateOnly(gender.created_at)}</td>
                  <td className="py-4 px-6">
                    <div className="flex items-center justify-center gap-2">
                      <button 
                        onClick={() => handleEdit(gender)}
                        className="text-blue-600 hover:text-blue-700 px-3 py-1 rounded hover:bg-blue-50 transition-colors"
                      >
                        แก้ไข
                      </button>
                      <button 
                        onClick={() => handleDelete(gender)}
                        className="text-red-600 hover:text-red-700 px-3 py-1 rounded hover:bg-red-50 transition-colors"
                      >
                        ลบ
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <GenderFormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        gender={selectedGender}
        mode={modalMode}
      />
    </div>
  );
}
