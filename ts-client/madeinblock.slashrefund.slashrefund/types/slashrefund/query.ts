/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";
import { Params } from "../slashrefund/params";
import { Deposit } from "../slashrefund/deposit";
import {
  PageRequest,
  PageResponse,
} from "../cosmos/base/query/v1beta1/pagination";
import { DepositPool } from "../slashrefund/deposit_pool";
import { UnbondingDeposit } from "../slashrefund/unbonding_deposit";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryGetDepositRequest {
  depositorAddress: string;
  validatorAddress: string;
}

export interface QueryGetDepositResponse {
  deposit: Deposit | undefined;
}

export interface QueryAllDepositRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllDepositResponse {
  deposit: Deposit[];
  pagination: PageResponse | undefined;
}

export interface QueryGetDepositPoolRequest {
  operatorAddress: string;
}

export interface QueryGetDepositPoolResponse {
  depositPool: DepositPool | undefined;
}

export interface QueryAllDepositPoolRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllDepositPoolResponse {
  depositPool: DepositPool[];
  pagination: PageResponse | undefined;
}

export interface QueryGetUnbondingDepositRequest {
  depositorAddress: string;
  validatorAddress: string;
}

export interface QueryGetUnbondingDepositResponse {
  unbondingDeposit: UnbondingDeposit | undefined;
}

export interface QueryAllUnbondingDepositRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllUnbondingDepositResponse {
  unbondingDeposit: UnbondingDeposit[];
  pagination: PageResponse | undefined;
}

const baseQueryParamsRequest: object = {};

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryParamsRequest } as QueryParamsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryParamsRequest {
    const message = { ...baseQueryParamsRequest } as QueryParamsRequest;
    return message;
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<QueryParamsRequest>): QueryParamsRequest {
    const message = { ...baseQueryParamsRequest } as QueryParamsRequest;
    return message;
  },
};

const baseQueryParamsResponse: object = {};

export const QueryParamsResponse = {
  encode(
    message: QueryParamsResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryParamsResponse } as QueryParamsResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryParamsResponse {
    const message = { ...baseQueryParamsResponse } as QueryParamsResponse;
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromJSON(object.params);
    } else {
      message.params = undefined;
    }
    return message;
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined &&
      (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryParamsResponse>): QueryParamsResponse {
    const message = { ...baseQueryParamsResponse } as QueryParamsResponse;
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    } else {
      message.params = undefined;
    }
    return message;
  },
};

const baseQueryGetDepositRequest: object = {
  depositorAddress: "",
  validatorAddress: "",
};

export const QueryGetDepositRequest = {
  encode(
    message: QueryGetDepositRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.depositorAddress !== "") {
      writer.uint32(10).string(message.depositorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetDepositRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryGetDepositRequest } as QueryGetDepositRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.depositorAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetDepositRequest {
    const message = { ...baseQueryGetDepositRequest } as QueryGetDepositRequest;
    if (
      object.depositorAddress !== undefined &&
      object.depositorAddress !== null
    ) {
      message.depositorAddress = String(object.depositorAddress);
    } else {
      message.depositorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = String(object.validatorAddress);
    } else {
      message.validatorAddress = "";
    }
    return message;
  },

  toJSON(message: QueryGetDepositRequest): unknown {
    const obj: any = {};
    message.depositorAddress !== undefined &&
      (obj.depositorAddress = message.depositorAddress);
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetDepositRequest>
  ): QueryGetDepositRequest {
    const message = { ...baseQueryGetDepositRequest } as QueryGetDepositRequest;
    if (
      object.depositorAddress !== undefined &&
      object.depositorAddress !== null
    ) {
      message.depositorAddress = object.depositorAddress;
    } else {
      message.depositorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = object.validatorAddress;
    } else {
      message.validatorAddress = "";
    }
    return message;
  },
};

const baseQueryGetDepositResponse: object = {};

export const QueryGetDepositResponse = {
  encode(
    message: QueryGetDepositResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.deposit !== undefined) {
      Deposit.encode(message.deposit, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetDepositResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetDepositResponse,
    } as QueryGetDepositResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.deposit = Deposit.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetDepositResponse {
    const message = {
      ...baseQueryGetDepositResponse,
    } as QueryGetDepositResponse;
    if (object.deposit !== undefined && object.deposit !== null) {
      message.deposit = Deposit.fromJSON(object.deposit);
    } else {
      message.deposit = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetDepositResponse): unknown {
    const obj: any = {};
    message.deposit !== undefined &&
      (obj.deposit = message.deposit
        ? Deposit.toJSON(message.deposit)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetDepositResponse>
  ): QueryGetDepositResponse {
    const message = {
      ...baseQueryGetDepositResponse,
    } as QueryGetDepositResponse;
    if (object.deposit !== undefined && object.deposit !== null) {
      message.deposit = Deposit.fromPartial(object.deposit);
    } else {
      message.deposit = undefined;
    }
    return message;
  },
};

const baseQueryAllDepositRequest: object = {};

export const QueryAllDepositRequest = {
  encode(
    message: QueryAllDepositRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryAllDepositRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryAllDepositRequest } as QueryAllDepositRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllDepositRequest {
    const message = { ...baseQueryAllDepositRequest } as QueryAllDepositRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllDepositRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllDepositRequest>
  ): QueryAllDepositRequest {
    const message = { ...baseQueryAllDepositRequest } as QueryAllDepositRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllDepositResponse: object = {};

export const QueryAllDepositResponse = {
  encode(
    message: QueryAllDepositResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.deposit) {
      Deposit.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryAllDepositResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllDepositResponse,
    } as QueryAllDepositResponse;
    message.deposit = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.deposit.push(Deposit.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllDepositResponse {
    const message = {
      ...baseQueryAllDepositResponse,
    } as QueryAllDepositResponse;
    message.deposit = [];
    if (object.deposit !== undefined && object.deposit !== null) {
      for (const e of object.deposit) {
        message.deposit.push(Deposit.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllDepositResponse): unknown {
    const obj: any = {};
    if (message.deposit) {
      obj.deposit = message.deposit.map((e) =>
        e ? Deposit.toJSON(e) : undefined
      );
    } else {
      obj.deposit = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllDepositResponse>
  ): QueryAllDepositResponse {
    const message = {
      ...baseQueryAllDepositResponse,
    } as QueryAllDepositResponse;
    message.deposit = [];
    if (object.deposit !== undefined && object.deposit !== null) {
      for (const e of object.deposit) {
        message.deposit.push(Deposit.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryGetDepositPoolRequest: object = { operatorAddress: "" };

export const QueryGetDepositPoolRequest = {
  encode(
    message: QueryGetDepositPoolRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.operatorAddress !== "") {
      writer.uint32(10).string(message.operatorAddress);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetDepositPoolRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetDepositPoolRequest,
    } as QueryGetDepositPoolRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.operatorAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetDepositPoolRequest {
    const message = {
      ...baseQueryGetDepositPoolRequest,
    } as QueryGetDepositPoolRequest;
    if (
      object.operatorAddress !== undefined &&
      object.operatorAddress !== null
    ) {
      message.operatorAddress = String(object.operatorAddress);
    } else {
      message.operatorAddress = "";
    }
    return message;
  },

  toJSON(message: QueryGetDepositPoolRequest): unknown {
    const obj: any = {};
    message.operatorAddress !== undefined &&
      (obj.operatorAddress = message.operatorAddress);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetDepositPoolRequest>
  ): QueryGetDepositPoolRequest {
    const message = {
      ...baseQueryGetDepositPoolRequest,
    } as QueryGetDepositPoolRequest;
    if (
      object.operatorAddress !== undefined &&
      object.operatorAddress !== null
    ) {
      message.operatorAddress = object.operatorAddress;
    } else {
      message.operatorAddress = "";
    }
    return message;
  },
};

const baseQueryGetDepositPoolResponse: object = {};

export const QueryGetDepositPoolResponse = {
  encode(
    message: QueryGetDepositPoolResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.depositPool !== undefined) {
      DepositPool.encode(
        message.depositPool,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetDepositPoolResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetDepositPoolResponse,
    } as QueryGetDepositPoolResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.depositPool = DepositPool.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetDepositPoolResponse {
    const message = {
      ...baseQueryGetDepositPoolResponse,
    } as QueryGetDepositPoolResponse;
    if (object.depositPool !== undefined && object.depositPool !== null) {
      message.depositPool = DepositPool.fromJSON(object.depositPool);
    } else {
      message.depositPool = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetDepositPoolResponse): unknown {
    const obj: any = {};
    message.depositPool !== undefined &&
      (obj.depositPool = message.depositPool
        ? DepositPool.toJSON(message.depositPool)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetDepositPoolResponse>
  ): QueryGetDepositPoolResponse {
    const message = {
      ...baseQueryGetDepositPoolResponse,
    } as QueryGetDepositPoolResponse;
    if (object.depositPool !== undefined && object.depositPool !== null) {
      message.depositPool = DepositPool.fromPartial(object.depositPool);
    } else {
      message.depositPool = undefined;
    }
    return message;
  },
};

const baseQueryAllDepositPoolRequest: object = {};

export const QueryAllDepositPoolRequest = {
  encode(
    message: QueryAllDepositPoolRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllDepositPoolRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllDepositPoolRequest,
    } as QueryAllDepositPoolRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllDepositPoolRequest {
    const message = {
      ...baseQueryAllDepositPoolRequest,
    } as QueryAllDepositPoolRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllDepositPoolRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllDepositPoolRequest>
  ): QueryAllDepositPoolRequest {
    const message = {
      ...baseQueryAllDepositPoolRequest,
    } as QueryAllDepositPoolRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllDepositPoolResponse: object = {};

export const QueryAllDepositPoolResponse = {
  encode(
    message: QueryAllDepositPoolResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.depositPool) {
      DepositPool.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllDepositPoolResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllDepositPoolResponse,
    } as QueryAllDepositPoolResponse;
    message.depositPool = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.depositPool.push(DepositPool.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllDepositPoolResponse {
    const message = {
      ...baseQueryAllDepositPoolResponse,
    } as QueryAllDepositPoolResponse;
    message.depositPool = [];
    if (object.depositPool !== undefined && object.depositPool !== null) {
      for (const e of object.depositPool) {
        message.depositPool.push(DepositPool.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllDepositPoolResponse): unknown {
    const obj: any = {};
    if (message.depositPool) {
      obj.depositPool = message.depositPool.map((e) =>
        e ? DepositPool.toJSON(e) : undefined
      );
    } else {
      obj.depositPool = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllDepositPoolResponse>
  ): QueryAllDepositPoolResponse {
    const message = {
      ...baseQueryAllDepositPoolResponse,
    } as QueryAllDepositPoolResponse;
    message.depositPool = [];
    if (object.depositPool !== undefined && object.depositPool !== null) {
      for (const e of object.depositPool) {
        message.depositPool.push(DepositPool.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryGetUnbondingDepositRequest: object = {
  depositorAddress: "",
  validatorAddress: "",
};

export const QueryGetUnbondingDepositRequest = {
  encode(
    message: QueryGetUnbondingDepositRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.depositorAddress !== "") {
      writer.uint32(10).string(message.depositorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetUnbondingDepositRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetUnbondingDepositRequest,
    } as QueryGetUnbondingDepositRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.depositorAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetUnbondingDepositRequest {
    const message = {
      ...baseQueryGetUnbondingDepositRequest,
    } as QueryGetUnbondingDepositRequest;
    if (
      object.depositorAddress !== undefined &&
      object.depositorAddress !== null
    ) {
      message.depositorAddress = String(object.depositorAddress);
    } else {
      message.depositorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = String(object.validatorAddress);
    } else {
      message.validatorAddress = "";
    }
    return message;
  },

  toJSON(message: QueryGetUnbondingDepositRequest): unknown {
    const obj: any = {};
    message.depositorAddress !== undefined &&
      (obj.depositorAddress = message.depositorAddress);
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetUnbondingDepositRequest>
  ): QueryGetUnbondingDepositRequest {
    const message = {
      ...baseQueryGetUnbondingDepositRequest,
    } as QueryGetUnbondingDepositRequest;
    if (
      object.depositorAddress !== undefined &&
      object.depositorAddress !== null
    ) {
      message.depositorAddress = object.depositorAddress;
    } else {
      message.depositorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = object.validatorAddress;
    } else {
      message.validatorAddress = "";
    }
    return message;
  },
};

const baseQueryGetUnbondingDepositResponse: object = {};

export const QueryGetUnbondingDepositResponse = {
  encode(
    message: QueryGetUnbondingDepositResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.unbondingDeposit !== undefined) {
      UnbondingDeposit.encode(
        message.unbondingDeposit,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetUnbondingDepositResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetUnbondingDepositResponse,
    } as QueryGetUnbondingDepositResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.unbondingDeposit = UnbondingDeposit.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetUnbondingDepositResponse {
    const message = {
      ...baseQueryGetUnbondingDepositResponse,
    } as QueryGetUnbondingDepositResponse;
    if (
      object.unbondingDeposit !== undefined &&
      object.unbondingDeposit !== null
    ) {
      message.unbondingDeposit = UnbondingDeposit.fromJSON(
        object.unbondingDeposit
      );
    } else {
      message.unbondingDeposit = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetUnbondingDepositResponse): unknown {
    const obj: any = {};
    message.unbondingDeposit !== undefined &&
      (obj.unbondingDeposit = message.unbondingDeposit
        ? UnbondingDeposit.toJSON(message.unbondingDeposit)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetUnbondingDepositResponse>
  ): QueryGetUnbondingDepositResponse {
    const message = {
      ...baseQueryGetUnbondingDepositResponse,
    } as QueryGetUnbondingDepositResponse;
    if (
      object.unbondingDeposit !== undefined &&
      object.unbondingDeposit !== null
    ) {
      message.unbondingDeposit = UnbondingDeposit.fromPartial(
        object.unbondingDeposit
      );
    } else {
      message.unbondingDeposit = undefined;
    }
    return message;
  },
};

const baseQueryAllUnbondingDepositRequest: object = {};

export const QueryAllUnbondingDepositRequest = {
  encode(
    message: QueryAllUnbondingDepositRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllUnbondingDepositRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllUnbondingDepositRequest,
    } as QueryAllUnbondingDepositRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllUnbondingDepositRequest {
    const message = {
      ...baseQueryAllUnbondingDepositRequest,
    } as QueryAllUnbondingDepositRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllUnbondingDepositRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllUnbondingDepositRequest>
  ): QueryAllUnbondingDepositRequest {
    const message = {
      ...baseQueryAllUnbondingDepositRequest,
    } as QueryAllUnbondingDepositRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllUnbondingDepositResponse: object = {};

export const QueryAllUnbondingDepositResponse = {
  encode(
    message: QueryAllUnbondingDepositResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.unbondingDeposit) {
      UnbondingDeposit.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllUnbondingDepositResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllUnbondingDepositResponse,
    } as QueryAllUnbondingDepositResponse;
    message.unbondingDeposit = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.unbondingDeposit.push(
            UnbondingDeposit.decode(reader, reader.uint32())
          );
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllUnbondingDepositResponse {
    const message = {
      ...baseQueryAllUnbondingDepositResponse,
    } as QueryAllUnbondingDepositResponse;
    message.unbondingDeposit = [];
    if (
      object.unbondingDeposit !== undefined &&
      object.unbondingDeposit !== null
    ) {
      for (const e of object.unbondingDeposit) {
        message.unbondingDeposit.push(UnbondingDeposit.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllUnbondingDepositResponse): unknown {
    const obj: any = {};
    if (message.unbondingDeposit) {
      obj.unbondingDeposit = message.unbondingDeposit.map((e) =>
        e ? UnbondingDeposit.toJSON(e) : undefined
      );
    } else {
      obj.unbondingDeposit = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllUnbondingDepositResponse>
  ): QueryAllUnbondingDepositResponse {
    const message = {
      ...baseQueryAllUnbondingDepositResponse,
    } as QueryAllUnbondingDepositResponse;
    message.unbondingDeposit = [];
    if (
      object.unbondingDeposit !== undefined &&
      object.unbondingDeposit !== null
    ) {
      for (const e of object.unbondingDeposit) {
        message.unbondingDeposit.push(UnbondingDeposit.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a Deposit by index. */
  Deposit(request: QueryGetDepositRequest): Promise<QueryGetDepositResponse>;
  /** Queries a list of Deposit items. */
  DepositAll(request: QueryAllDepositRequest): Promise<QueryAllDepositResponse>;
  /** Queries a DepositPool by index. */
  DepositPool(
    request: QueryGetDepositPoolRequest
  ): Promise<QueryGetDepositPoolResponse>;
  /** Queries a list of DepositPool items. */
  DepositPoolAll(
    request: QueryAllDepositPoolRequest
  ): Promise<QueryAllDepositPoolResponse>;
  /** Queries a UnbondingDeposit by index. */
  UnbondingDeposit(
    request: QueryGetUnbondingDepositRequest
  ): Promise<QueryGetUnbondingDepositResponse>;
  /** Queries a list of UnbondingDeposit items. */
  UnbondingDepositAll(
    request: QueryAllUnbondingDepositRequest
  ): Promise<QueryAllUnbondingDepositResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "madeinblock.slashrefund.slashrefund.Query",
      "Params",
      data
    );
    return promise.then((data) => QueryParamsResponse.decode(new Reader(data)));
  }

  Deposit(request: QueryGetDepositRequest): Promise<QueryGetDepositResponse> {
    const data = QueryGetDepositRequest.encode(request).finish();
    const promise = this.rpc.request(
      "madeinblock.slashrefund.slashrefund.Query",
      "Deposit",
      data
    );
    return promise.then((data) =>
      QueryGetDepositResponse.decode(new Reader(data))
    );
  }

  DepositAll(
    request: QueryAllDepositRequest
  ): Promise<QueryAllDepositResponse> {
    const data = QueryAllDepositRequest.encode(request).finish();
    const promise = this.rpc.request(
      "madeinblock.slashrefund.slashrefund.Query",
      "DepositAll",
      data
    );
    return promise.then((data) =>
      QueryAllDepositResponse.decode(new Reader(data))
    );
  }

  DepositPool(
    request: QueryGetDepositPoolRequest
  ): Promise<QueryGetDepositPoolResponse> {
    const data = QueryGetDepositPoolRequest.encode(request).finish();
    const promise = this.rpc.request(
      "madeinblock.slashrefund.slashrefund.Query",
      "DepositPool",
      data
    );
    return promise.then((data) =>
      QueryGetDepositPoolResponse.decode(new Reader(data))
    );
  }

  DepositPoolAll(
    request: QueryAllDepositPoolRequest
  ): Promise<QueryAllDepositPoolResponse> {
    const data = QueryAllDepositPoolRequest.encode(request).finish();
    const promise = this.rpc.request(
      "madeinblock.slashrefund.slashrefund.Query",
      "DepositPoolAll",
      data
    );
    return promise.then((data) =>
      QueryAllDepositPoolResponse.decode(new Reader(data))
    );
  }

  UnbondingDeposit(
    request: QueryGetUnbondingDepositRequest
  ): Promise<QueryGetUnbondingDepositResponse> {
    const data = QueryGetUnbondingDepositRequest.encode(request).finish();
    const promise = this.rpc.request(
      "madeinblock.slashrefund.slashrefund.Query",
      "UnbondingDeposit",
      data
    );
    return promise.then((data) =>
      QueryGetUnbondingDepositResponse.decode(new Reader(data))
    );
  }

  UnbondingDepositAll(
    request: QueryAllUnbondingDepositRequest
  ): Promise<QueryAllUnbondingDepositResponse> {
    const data = QueryAllUnbondingDepositRequest.encode(request).finish();
    const promise = this.rpc.request(
      "madeinblock.slashrefund.slashrefund.Query",
      "UnbondingDepositAll",
      data
    );
    return promise.then((data) =>
      QueryAllUnbondingDepositResponse.decode(new Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;
