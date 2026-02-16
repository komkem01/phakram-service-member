import { apiClient } from './client';

export interface Province {
  id: string;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateProvinceDto {
  name: string;
  is_active: boolean;
}

export interface UpdateProvinceDto {
  name: string;
  is_active: boolean;
}

export interface ListProvincesParams {
  page?: number;
  size?: number;
}

interface BackendPaginate {
  page?: number;
  size?: number;
  total?: number;
  Page?: number;
  Size?: number;
  Total?: number;
}

export interface ListProvincesResponse {
  data: Province[];
  paginate: {
    page: number;
    size: number;
    total: number;
  };
}

export const provincesApi = {
  async list(params?: ListProvincesParams): Promise<ListProvincesResponse> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.size) queryParams.append('size', params.size.toString());

    const query = queryParams.toString();
    const endpoint = `/system/provinces/${query ? `?${query}` : ''}`;

    const response = await apiClient.get<Province[]>(endpoint);
    const rawPaginate = (response.paginate ?? {}) as BackendPaginate;

    return {
      data: response.data,
      paginate: {
        page: rawPaginate.page ?? rawPaginate.Page ?? 1,
        size: rawPaginate.size ?? rawPaginate.Size ?? params?.size ?? 10,
        total: rawPaginate.total ?? rawPaginate.Total ?? response.data.length,
      },
    };
  },

  async getById(id: string): Promise<Province> {
    const response = await apiClient.get<Province>(`/system/provinces/${id}`);
    return response.data;
  },

  async create(data: CreateProvinceDto): Promise<void> {
    await apiClient.post('/system/provinces/', data);
  },

  async update(id: string, data: UpdateProvinceDto): Promise<void> {
    await apiClient.patch(`/system/provinces/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/system/provinces/${id}`);
  },
};
