'use client';

import { useState, useEffect, useCallback } from 'react';
import Loading from '@/components/admin/Loading';
import BankFormModal from '@/components/admin/BankFormModal';
import { useModal } from '@/contexts/ModalContext';
import { banksApi, Bank, BankCreateInput, BankUpdateInput } from '@/lib/api/banks';
import { formatDateOnly } from '@/lib/utils/date';

export default function BanksPage() {
  const [banks, setBanks] = useState<Bank[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [selectedBank, setSelectedBank] = useState<Bank | null>(null);
  const { showSuccess, showWarning, showError } = useModal();

  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  const fetchBanks = useCallback(async (currentPage: number, currentSize: number) => {
    try {
      setIsLoading(true);
      const response = await banksApi.list({
        page: currentPage,
        size: currentSize,
      });

      setBanks(response.data);
      setTotal(response.paginate.total);
      setPage(response.paginate.page);
    } catch (error) {
      console.error('Failed to fetch banks:', error);
      showError('ไม่สามารถโหลดข้อมูลธนาคารได้', 'เกิดข้อผิดพลาด');
    } finally {
      setIsLoading(false);
    }
  }, [showError]);

  useEffect(() => {
    fetchBanks(page, pageSize);
  }, [fetchBanks, page, pageSize]);

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
    setSelectedBank(null);
    setIsModalOpen(true);
  };

  const handleEdit = (bank: Bank) => {
    setModalMode('edit');
    setSelectedBank(bank);
    setIsModalOpen(true);
  };

  const handleDelete = (bank: Bank) => {
    showWarning(
      `คุณต้องการลบธนาคาร "${bank.name_th}" ใช่หรือไม่?`,
      async () => {
        try {
          await banksApi.delete(bank.id);
          showSuccess(`ลบธนาคาร "${bank.name_th}" สำเร็จ`);
          fetchBanks(page, pageSize);
        } catch (error) {
          console.error('Failed to delete bank:', error);
          showError(`ไม่สามารถลบธนาคาร "${bank.name_th}" ได้`, 'เกิดข้อผิดพลาด');
        }
      },
      'ยืนยันการลบ'
    );
  };

  const handleSubmit = async (data: BankCreateInput) => {
    try {
      if (modalMode === 'create') {
        await banksApi.create(data);
        showSuccess(`เพิ่มธนาคาร "${data.name_th}" สำเร็จ`);
      } else if (modalMode === 'edit' && selectedBank) {
        const updateData: BankUpdateInput = {
          name_th: data.name_th,
          name_abb_th: data.name_abb_th,
          name_en: data.name_en,
          name_abb_en: data.name_abb_en,
          is_active: data.is_active,
        };
        await banksApi.update(selectedBank.id, updateData);
        showSuccess(`แก้ไขธนาคาร "${data.name_th}" สำเร็จ`);
      }
      fetchBanks(page, pageSize);
    } catch (error) {
      console.error('Failed to save bank:', error);
      showError(
        `ไม่สามารถ${modalMode === 'create' ? 'เพิ่ม' : 'แก้ไข'}ธนาคารได้`,
        'เกิดข้อผิดพลาด'
      );
    }
  };

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการธนาคาร</h1>
        <button
          onClick={handleCreate}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          เพิ่มธนาคาร
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="bg-gray-50 border-b border-gray-200">
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ชื่อธนาคาร (ไทย)
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ชื่อย่อ (ไทย)
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ชื่อธนาคาร (อังกฤษ)
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ชื่อย่อ (อังกฤษ)
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
              {banks.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-6 py-12 text-center text-gray-500">
                    ไม่พบข้อมูลธนาคาร
                  </td>
                </tr>
              ) : (
                banks.map((bank) => (
                  <tr key={bank.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {bank.name_th}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {bank.name_abb_th}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {bank.name_en}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {bank.name_abb_en}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-center">
                      <span
                        className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                          bank.is_active
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {bank.is_active ? 'ใช้งาน' : 'ไม่ใช้งาน'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {formatDateOnly(bank.created_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-center text-sm font-medium">
                      <button
                        onClick={() => handleEdit(bank)}
                        className="text-blue-600 hover:text-blue-900 mr-3"
                      >
                        แก้ไข
                      </button>
                      <button
                        onClick={() => handleDelete(bank)}
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

      <BankFormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        bank={selectedBank}
        mode={modalMode}
      />
    </div>
  );
}
