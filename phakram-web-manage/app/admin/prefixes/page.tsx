"use client";

import { useState, useEffect, useCallback } from "react";
import Loading from "@/components/admin/Loading";
import PrefixFormModal from "@/components/admin/PrefixFormModal";
import { useModal } from "@/contexts/ModalContext";
import { prefixesApi, Prefix, CreatePrefixDto, UpdatePrefixDto } from "@/lib/api/prefixes";
import { gendersApi } from "@/lib/api/genders";
import { formatDateOnly } from "@/lib/utils/date";

export default function PrefixesPage() {
  const [prefixes, setPrefixes] = useState<Prefix[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [selectedPrefix, setSelectedPrefix] = useState<Prefix | null>(null);
  const [genderNames, setGenderNames] = useState<Record<string, string>>({});
  const { showSuccess, showWarning, showError } = useModal();

  // Fetch prefixes from API
  const fetchPrefixes = useCallback(async () => {
    try {
      setIsLoading(true);
      const [prefixResponse, genderResponse] = await Promise.all([
        prefixesApi.list(),
        gendersApi.list(),
      ]);
      setPrefixes(prefixResponse.data);
      
      // Create gender name lookup
      const genderMap: Record<string, string> = {};
      genderResponse.data.forEach((gender) => {
        genderMap[gender.id] = gender.name_th;
      });
      setGenderNames(genderMap);
    } catch (error) {
      console.error('Failed to fetch prefixes:', error);
      showError('ไม่สามารถโหลดข้อมูลคำนำหน้าได้', 'เกิดข้อผิดพลาด');
    } finally {
      setIsLoading(false);
    }
  }, [showError]);

  useEffect(() => {
    fetchPrefixes();
  }, [fetchPrefixes]);

  const handleCreate = () => {
    setModalMode('create');
    setSelectedPrefix(null);
    setIsModalOpen(true);
  };

  const handleEdit = (prefix: Prefix) => {
    setModalMode('edit');
    setSelectedPrefix(prefix);
    setIsModalOpen(true);
  };

  const handleDelete = (prefix: Prefix) => {
    showWarning(
      `คุณต้องการลบคำนำหน้า "${prefix.name_th}" ใช่หรือไม่?`,
      async () => {
        try {
          await prefixesApi.delete(prefix.id);
          showSuccess(`ลบคำนำหน้า "${prefix.name_th}" สำเร็จ`);
          fetchPrefixes();
        } catch (error) {
          console.error('Failed to delete prefix:', error);
          showError(`ไม่สามารถลบคำนำหน้า "${prefix.name_th}" ได้`, 'เกิดข้อผิดพลาด');
        }
      },
      'ยืนยันการลบ'
    );
  };

  const handleSubmit = async (data: Omit<Prefix, 'id' | 'created_at'>) => {
    try {
      if (modalMode === 'create') {
        const createData: CreatePrefixDto = {
          name_th: data.name_th,
          name_en: data.name_en,
          gender_id: data.gender_id,
          is_active: data.is_active,
        };
        await prefixesApi.create(createData);
        showSuccess(`เพิ่มคำนำหน้า "${data.name_th}" สำเร็จ`);
      } else if (modalMode === 'edit' && selectedPrefix) {
        const updateData: UpdatePrefixDto = {
          name_th: data.name_th,
          name_en: data.name_en,
          gender_id: data.gender_id,
          is_active: data.is_active,
        };
        await prefixesApi.update(selectedPrefix.id, updateData);
        showSuccess(`แก้ไขคำนำหน้า "${data.name_th}" สำเร็จ`);
      }
      fetchPrefixes();
    } catch (error) {
      console.error('Failed to save prefix:', error);
      showError(
        `ไม่สามารถ${modalMode === 'create' ? 'เพิ่ม' : 'แก้ไข'}คำนำหน้าได้`,
        'เกิดข้อผิดพลาด'
      );
    }
  };

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการคำนำหน้า</h1>
        <button 
          onClick={handleCreate}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          + เพิ่มคำนำหน้า
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50">
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ลำดับ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">คำนำหน้า (ไทย)</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">คำนำหน้า (อังกฤษ)</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">เพศ</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">สถานะ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">วันที่สร้าง</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              {prefixes.map((prefix, index) => (
                <tr key={prefix.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-4 px-6">{index + 1}</td>
                  <td className="py-4 px-6 font-medium">{prefix.name_th}</td>
                  <td className="py-4 px-6 text-gray-600">{prefix.name_en}</td>
                  <td className="py-4 px-6 text-gray-600">
                    {genderNames[prefix.gender_id] || '-'}
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex justify-center">
                      <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                        prefix.is_active 
                          ? 'bg-green-100 text-green-700' 
                          : 'bg-gray-100 text-gray-600'
                      }`}>
                        {prefix.is_active ? 'เปิดใช้งาน' : 'ปิดใช้งาน'}
                      </span>
                    </div>
                  </td>
                  <td className="py-4 px-6 text-gray-600">{formatDateOnly(prefix.created_at)}</td>
                  <td className="py-4 px-6">
                    <div className="flex items-center justify-center gap-2">
                      <button 
                        onClick={() => handleEdit(prefix)}
                        className="text-blue-600 hover:text-blue-700 px-3 py-1 rounded hover:bg-blue-50 transition-colors"
                      >
                        แก้ไข
                      </button>
                      <button 
                        onClick={() => handleDelete(prefix)}
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

      <PrefixFormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        prefix={selectedPrefix}
        mode={modalMode}
      />
    </div>
  );
}
