'use client';

import { useState, useEffect, useCallback } from 'react';
import Loading from '@/components/admin/Loading';
import TierFormModal from '@/components/admin/TierFormModal';
import { useModal } from '@/contexts/ModalContext';
import { tiersApi, Tier, TierCreateInput, TierUpdateInput } from '@/lib/api/tiers';
import { formatDateOnly } from '@/lib/utils/date';

export default function TiersPage() {
  const [tiers, setTiers] = useState<Tier[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [selectedTier, setSelectedTier] = useState<Tier | null>(null);
  const { showSuccess, showWarning, showError } = useModal();

  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  const fetchTiers = useCallback(async (currentPage: number, currentSize: number) => {
    try {
      setIsLoading(true);
      const response = await tiersApi.list({
        page: currentPage,
        size: currentSize,
      });

      setTiers(response.data);
      setTotal(response.paginate.total);
      setPage(response.paginate.page);
    } catch (error) {
      console.error('Failed to fetch tiers:', error);
      showError('ไม่สามารถโหลดข้อมูลระดับสมาชิกได้', 'เกิดข้อผิดพลาด');
    } finally {
      setIsLoading(false);
    }
  }, [showError]);

  useEffect(() => {
    fetchTiers(page, pageSize);
  }, [fetchTiers, page, pageSize]);

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
    setSelectedTier(null);
    setIsModalOpen(true);
  };

  const handleEdit = (tier: Tier) => {
    setModalMode('edit');
    setSelectedTier(tier);
    setIsModalOpen(true);
  };

  const handleDelete = (tier: Tier) => {
    showWarning(
      `คุณต้องการลบระดับสมาชิก "${tier.name_th}" ใช่หรือไม่?`,
      async () => {
        try {
          await tiersApi.delete(tier.id);
          showSuccess(`ลบระดับสมาชิก "${tier.name_th}" สำเร็จ`);
          fetchTiers(page, pageSize);
        } catch (error) {
          console.error('Failed to delete tier:', error);
          showError(`ไม่สามารถลบระดับสมาชิก "${tier.name_th}" ได้`, 'เกิดข้อผิดพลาด');
        }
      },
      'ยืนยันการลบ'
    );
  };

  const handleSubmit = async (data: TierCreateInput) => {
    try {
      if (modalMode === 'create') {
        await tiersApi.create(data);
        showSuccess(`เพิ่มระดับสมาชิก "${data.name_th}" สำเร็จ`);
      } else if (modalMode === 'edit' && selectedTier) {
        const updateData: TierUpdateInput = {
          name_th: data.name_th,
          name_en: data.name_en,
          min_spending: data.min_spending,
          discount_rate: data.discount_rate,
          is_active: data.is_active,
        };
        await tiersApi.update(selectedTier.id, updateData);
        showSuccess(`แก้ไขระดับสมาชิก "${data.name_th}" สำเร็จ`);
      }
      fetchTiers(page, pageSize);
    } catch (error) {
      console.error('Failed to save tier:', error);
      showError(
        `ไม่สามารถ${modalMode === 'create' ? 'เพิ่ม' : 'แก้ไข'}ระดับสมาชิกได้`,
        'เกิดข้อผิดพลาด'
      );
    }
  };

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการระดับสมาชิก</h1>
        <button
          onClick={handleCreate}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          เพิ่มระดับสมาชิก
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="bg-gray-50 border-b border-gray-200">
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ชื่อ (ไทย)
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ชื่อ (อังกฤษ)
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ยอดใช้จ่ายขั้นต่ำ
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ส่วนลด (%)
                </th>
                <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                  สถานะ
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  สร้างเมื่อ
                </th>
                <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                  จัดการ
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {tiers.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-6 py-12 text-center text-gray-500">
                    ไม่พบข้อมูลระดับสมาชิก
                  </td>
                </tr>
              ) : (
                tiers.map((tier) => (
                  <tr key={tier.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {tier.name_th}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {tier.name_en}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700 text-right">
                      {tier.min_spending.toLocaleString('th-TH', {
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                      })}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700 text-right">
                      {tier.discount_rate.toFixed(2)}%
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-center">
                      <span
                        className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                          tier.is_active
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {tier.is_active ? 'ใช้งาน' : 'ไม่ใช้งาน'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {formatDateOnly(tier.created_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-center text-sm font-medium">
                      <button
                        onClick={() => handleEdit(tier)}
                        className="text-blue-600 hover:text-blue-900 mr-3"
                      >
                        แก้ไข
                      </button>
                      <button
                        onClick={() => handleDelete(tier)}
                        className="text-red-600 hover:text-red-900"
                      >
                        ลบ
                      </button>
                    </td>
                  </tr>
                ))
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

      <TierFormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        tier={selectedTier}
        mode={modalMode}
      />
    </div>
  );
}
