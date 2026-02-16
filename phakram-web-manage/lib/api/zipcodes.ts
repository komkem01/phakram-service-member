import { apiClient } from './client';

export interface Zipcode {
  id: string;
  sub_districts_id: string;
  sub_district_id?: string;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateZipcodeDto {
  sub_districts_id: string;
  name: string;
  is_active: boolean;
}

export interface UpdateZipcodeDto {
  sub_districts_id: string;
  name: string;
  is_active: boolean;
}

export interface ListZipcodesParams {
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

export interface ListZipcodesResponse {
  data: Zipcode[];
  paginate: {
    page: number;
    size: number;
    total: number;
  };
}

export const zipcodesApi = {
  async list(params?: ListZipcodesParams): Promise<ListZipcodesResponse> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.size) queryParams.append('size', params.size.toString());

    const query = queryParams.toString();
    const endpoint = `/system/zipcodes/${query ? `?${query}` : ''}`;

    const response = await apiClient.get<Zipcode[]>(endpoint);
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

  async getById(id: string): Promise<Zipcode> {
    const response = await apiClient.get<Zipcode>(`/system/zipcodes/${id}`);
    return response.data;
  },

  async create(data: CreateZipcodeDto): Promise<void> {
    await apiClient.post('/system/zipcodes/', data);
  },

  async update(id: string, data: UpdateZipcodeDto): Promise<void> {
    await apiClient.patch(`/system/zipcodes/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/system/zipcodes/${id}`);
  },
};
