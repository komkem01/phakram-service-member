"use client";

import { useState, useEffect, useCallback } from "react";
import Loading from "@/components/admin/Loading";
import ProvinceFormModal from "@/components/admin/ProvinceFormModal";
import { useModal } from "@/contexts/ModalContext";
import { provincesApi, Province, CreateProvinceDto, UpdateProvinceDto } from "@/lib/api/provinces";
import { formatDateOnly } from "@/lib/utils/date";

export default function ProvincesPage() {
  const [provinces, setProvinces] = useState<Province[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [selectedProvince, setSelectedProvince] = useState<Province | null>(null);
  const { showSuccess, showWarning, showError } = useModal();

  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  const fetchProvinces = useCallback(async (currentPage: number, currentSize: number) => {
    try {
      setIsLoading(true);
      const response = await provincesApi.list({
        page: currentPage,
        size: currentSize,
      });

      setProvinces(response.data);
      setTotal(response.paginate.total);
      setPage(response.paginate.page);
    } catch (error) {
      console.error("Failed to fetch provinces:", error);
      showError("ไม่สามารถโหลดข้อมูลจังหวัดได้", "เกิดข้อผิดพลาด");
    } finally {
      setIsLoading(false);
    }
  }, [showError]);

  useEffect(() => {
    fetchProvinces(page, pageSize);
  }, [fetchProvinces, page, pageSize]);

  const handlePageChange = (newPage: number) => {
    if (newPage < 1 || newPage > totalPages) return;
    setPage(newPage);
  };

  const handlePageSizeChange = (value: number) => {
    setPageSize(value);
    setPage(1);
  };

  const handleCreate = () => {
    setModalMode('create');
    setSelectedProvince(null);
    setIsModalOpen(true);
  };

  const handleEdit = (province: Province) => {
    setModalMode('edit');
    setSelectedProvince(province);
    setIsModalOpen(true);
  };

  const handleDelete = (province: Province) => {
    showWarning(
      `คุณต้องการลบจังหวัด "${province.name}" ใช่หรือไม่?`,
      async () => {
        try {
          await provincesApi.delete(province.id);
          showSuccess(`ลบจังหวัด "${province.name}" สำเร็จ`);
          fetchProvinces(page, pageSize);
        } catch (error) {
          console.error('Failed to delete province:', error);
          showError(`ไม่สามารถลบจังหวัด "${province.name}" ได้`, 'เกิดข้อผิดพลาด');
        }
      },
      'ยืนยันการลบ'
    );
  };

  const handleSubmit = async (data: Omit<Province, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      if (modalMode === 'create') {
        const createData: CreateProvinceDto = {
          name: data.name,
          is_active: data.is_active,
        };
        await provincesApi.create(createData);
        showSuccess(`เพิ่มจังหวัด "${data.name}" สำเร็จ`);
      } else if (modalMode === 'edit' && selectedProvince) {
        const updateData: UpdateProvinceDto = {
          name: data.name,
          is_active: data.is_active,
        };
        await provincesApi.update(selectedProvince.id, updateData);
        showSuccess(`แก้ไขจังหวัด "${data.name}" สำเร็จ`);
      }
      fetchProvinces(page, pageSize);
    } catch (error) {
      console.error('Failed to save province:', error);
      showError(
        `ไม่สามารถ${modalMode === 'create' ? 'เพิ่ม' : 'แก้ไข'}จังหวัดได้`,
        'เกิดข้อผิดพลาด'
      );
    }
  };

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการจังหวัด</h1>
        <button 
          onClick={handleCreate}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          + เพิ่มจังหวัด
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50">
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ลำดับ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ชื่อจังหวัด</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">สถานะ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">วันที่สร้าง</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              {provinces.length > 0 ? (
                provinces.map((province, index) => (
                  <tr key={province.id} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="py-4 px-6">{(page - 1) * pageSize + index + 1}</td>
                    <td className="py-4 px-6 font-medium">{province.name}</td>
                    <td className="py-4 px-6">
                      <div className="flex justify-center">
                        <span
                          className={`px-3 py-1 rounded-full text-xs font-semibold ${
                            province.is_active
                              ? "bg-green-100 text-green-700"
                              : "bg-gray-100 text-gray-600"
                          }`}
                        >
                          {province.is_active ? "เปิดใช้งาน" : "ปิดใช้งาน"}
                        </span>
                      </div>
                    </td>
                    <td className="py-4 px-6 text-gray-600">{formatDateOnly(province.created_at)}</td>
                    <td className="py-4 px-6">
                      <div className="flex items-center justify-center gap-2">
                        <button 
                          onClick={() => handleEdit(province)}
                          className="text-blue-600 hover:text-blue-700 px-3 py-1 rounded hover:bg-blue-50 transition-colors"
                        >
                          แก้ไข
                        </button>
                        <button 
                          onClick={() => handleDelete(province)}
                          className="text-red-600 hover:text-red-700 px-3 py-1 rounded hover:bg-red-50 transition-colors"
                        >
                          ลบ
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={5} className="py-8 text-center text-gray-500">
                    ไม่พบข้อมูลจังหวัด
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between p-4 border-t border-gray-200">
          <div className="flex items-center gap-2 text-sm text-gray-600">
            <span>แสดง</span>
            <select
              value={pageSize}
              onChange={(e) => handlePageSizeChange(Number(e.target.value))}
              className="border border-gray-300 rounded-md px-2 py-1 bg-white"
            >
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
              <option value={100}>100</option>
            </select>
            <span>รายการต่อหน้า</span>
          </div>

          <div className="text-sm text-gray-600">
            หน้า {page} / {totalPages} (ทั้งหมด {total} รายการ)
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => handlePageChange(page - 1)}
              disabled={page <= 1}
              className="px-3 py-1 rounded-md border border-gray-300 text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              ก่อนหน้า
            </button>
            <button
              onClick={() => handlePageChange(page + 1)}
              disabled={page >= totalPages}
              className="px-3 py-1 rounded-md border border-gray-300 text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              ถัดไป
            </button>
          </div>
        </div>
      </div>

      <ProvinceFormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        province={selectedProvince}
        mode={modalMode}
      />
    </div>
  );
}
