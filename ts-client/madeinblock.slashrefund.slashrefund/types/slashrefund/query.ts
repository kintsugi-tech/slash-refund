/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../cosmos/base/query/v1beta1/pagination";
import { Deposit } from "./deposit";
import { DepositPool } from "./deposit_pool";
import { Params } from "./params";
import { Refund } from "./refund";
import { RefundPool } from "./refund_pool";
import { UnbondingDeposit } from "./unbonding_deposit";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

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

export interface QueryGetRefundPoolRequest {
  operatorAddress: string;
}

export interface QueryGetRefundPoolResponse {
  refundPool: RefundPool | undefined;
}

export interface QueryAllRefundPoolRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllRefundPoolResponse {
  refundPool: RefundPool[];
  pagination: PageResponse | undefined;
}

export interface QueryGetRefundRequest {
  delegator: string;
  validator: string;
}

export interface QueryGetRefundResponse {
  refund: Refund | undefined;
}

export interface QueryAllRefundRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllRefundResponse {
  refund: Refund[];
  pagination: PageResponse | undefined;
}

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
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
    return {};
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsRequest>, I>>(_: I): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return { params: undefined };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
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
    return { params: isSet(object.params) ? Params.fromJSON(object.params) : undefined };
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsResponse>, I>>(object: I): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    return message;
  },
};

function createBaseQueryGetDepositRequest(): QueryGetDepositRequest {
  return { depositorAddress: "", validatorAddress: "" };
}

export const QueryGetDepositRequest = {
  encode(message: QueryGetDepositRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.depositorAddress !== "") {
      writer.uint32(10).string(message.depositorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetDepositRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetDepositRequest();
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
    return {
      depositorAddress: isSet(object.depositorAddress) ? String(object.depositorAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
    };
  },

  toJSON(message: QueryGetDepositRequest): unknown {
    const obj: any = {};
    message.depositorAddress !== undefined && (obj.depositorAddress = message.depositorAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetDepositRequest>, I>>(object: I): QueryGetDepositRequest {
    const message = createBaseQueryGetDepositRequest();
    message.depositorAddress = object.depositorAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    return message;
  },
};

function createBaseQueryGetDepositResponse(): QueryGetDepositResponse {
  return { deposit: undefined };
}

export const QueryGetDepositResponse = {
  encode(message: QueryGetDepositResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.deposit !== undefined) {
      Deposit.encode(message.deposit, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetDepositResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetDepositResponse();
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
    return { deposit: isSet(object.deposit) ? Deposit.fromJSON(object.deposit) : undefined };
  },

  toJSON(message: QueryGetDepositResponse): unknown {
    const obj: any = {};
    message.deposit !== undefined && (obj.deposit = message.deposit ? Deposit.toJSON(message.deposit) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetDepositResponse>, I>>(object: I): QueryGetDepositResponse {
    const message = createBaseQueryGetDepositResponse();
    message.deposit = (object.deposit !== undefined && object.deposit !== null)
      ? Deposit.fromPartial(object.deposit)
      : undefined;
    return message;
  },
};

function createBaseQueryAllDepositRequest(): QueryAllDepositRequest {
  return { pagination: undefined };
}

export const QueryAllDepositRequest = {
  encode(message: QueryAllDepositRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDepositRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllDepositRequest();
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
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllDepositRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllDepositRequest>, I>>(object: I): QueryAllDepositRequest {
    const message = createBaseQueryAllDepositRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllDepositResponse(): QueryAllDepositResponse {
  return { deposit: [], pagination: undefined };
}

export const QueryAllDepositResponse = {
  encode(message: QueryAllDepositResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.deposit) {
      Deposit.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDepositResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllDepositResponse();
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
    return {
      deposit: Array.isArray(object?.deposit) ? object.deposit.map((e: any) => Deposit.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllDepositResponse): unknown {
    const obj: any = {};
    if (message.deposit) {
      obj.deposit = message.deposit.map((e) => e ? Deposit.toJSON(e) : undefined);
    } else {
      obj.deposit = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllDepositResponse>, I>>(object: I): QueryAllDepositResponse {
    const message = createBaseQueryAllDepositResponse();
    message.deposit = object.deposit?.map((e) => Deposit.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetDepositPoolRequest(): QueryGetDepositPoolRequest {
  return { operatorAddress: "" };
}

export const QueryGetDepositPoolRequest = {
  encode(message: QueryGetDepositPoolRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.operatorAddress !== "") {
      writer.uint32(10).string(message.operatorAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetDepositPoolRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetDepositPoolRequest();
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
    return { operatorAddress: isSet(object.operatorAddress) ? String(object.operatorAddress) : "" };
  },

  toJSON(message: QueryGetDepositPoolRequest): unknown {
    const obj: any = {};
    message.operatorAddress !== undefined && (obj.operatorAddress = message.operatorAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetDepositPoolRequest>, I>>(object: I): QueryGetDepositPoolRequest {
    const message = createBaseQueryGetDepositPoolRequest();
    message.operatorAddress = object.operatorAddress ?? "";
    return message;
  },
};

function createBaseQueryGetDepositPoolResponse(): QueryGetDepositPoolResponse {
  return { depositPool: undefined };
}

export const QueryGetDepositPoolResponse = {
  encode(message: QueryGetDepositPoolResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.depositPool !== undefined) {
      DepositPool.encode(message.depositPool, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetDepositPoolResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetDepositPoolResponse();
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
    return { depositPool: isSet(object.depositPool) ? DepositPool.fromJSON(object.depositPool) : undefined };
  },

  toJSON(message: QueryGetDepositPoolResponse): unknown {
    const obj: any = {};
    message.depositPool !== undefined
      && (obj.depositPool = message.depositPool ? DepositPool.toJSON(message.depositPool) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetDepositPoolResponse>, I>>(object: I): QueryGetDepositPoolResponse {
    const message = createBaseQueryGetDepositPoolResponse();
    message.depositPool = (object.depositPool !== undefined && object.depositPool !== null)
      ? DepositPool.fromPartial(object.depositPool)
      : undefined;
    return message;
  },
};

function createBaseQueryAllDepositPoolRequest(): QueryAllDepositPoolRequest {
  return { pagination: undefined };
}

export const QueryAllDepositPoolRequest = {
  encode(message: QueryAllDepositPoolRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDepositPoolRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllDepositPoolRequest();
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
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllDepositPoolRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllDepositPoolRequest>, I>>(object: I): QueryAllDepositPoolRequest {
    const message = createBaseQueryAllDepositPoolRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllDepositPoolResponse(): QueryAllDepositPoolResponse {
  return { depositPool: [], pagination: undefined };
}

export const QueryAllDepositPoolResponse = {
  encode(message: QueryAllDepositPoolResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.depositPool) {
      DepositPool.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDepositPoolResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllDepositPoolResponse();
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
    return {
      depositPool: Array.isArray(object?.depositPool)
        ? object.depositPool.map((e: any) => DepositPool.fromJSON(e))
        : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllDepositPoolResponse): unknown {
    const obj: any = {};
    if (message.depositPool) {
      obj.depositPool = message.depositPool.map((e) => e ? DepositPool.toJSON(e) : undefined);
    } else {
      obj.depositPool = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllDepositPoolResponse>, I>>(object: I): QueryAllDepositPoolResponse {
    const message = createBaseQueryAllDepositPoolResponse();
    message.depositPool = object.depositPool?.map((e) => DepositPool.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetUnbondingDepositRequest(): QueryGetUnbondingDepositRequest {
  return { depositorAddress: "", validatorAddress: "" };
}

export const QueryGetUnbondingDepositRequest = {
  encode(message: QueryGetUnbondingDepositRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.depositorAddress !== "") {
      writer.uint32(10).string(message.depositorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetUnbondingDepositRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetUnbondingDepositRequest();
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
    return {
      depositorAddress: isSet(object.depositorAddress) ? String(object.depositorAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
    };
  },

  toJSON(message: QueryGetUnbondingDepositRequest): unknown {
    const obj: any = {};
    message.depositorAddress !== undefined && (obj.depositorAddress = message.depositorAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetUnbondingDepositRequest>, I>>(
    object: I,
  ): QueryGetUnbondingDepositRequest {
    const message = createBaseQueryGetUnbondingDepositRequest();
    message.depositorAddress = object.depositorAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    return message;
  },
};

function createBaseQueryGetUnbondingDepositResponse(): QueryGetUnbondingDepositResponse {
  return { unbondingDeposit: undefined };
}

export const QueryGetUnbondingDepositResponse = {
  encode(message: QueryGetUnbondingDepositResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.unbondingDeposit !== undefined) {
      UnbondingDeposit.encode(message.unbondingDeposit, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetUnbondingDepositResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetUnbondingDepositResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.unbondingDeposit = UnbondingDeposit.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetUnbondingDepositResponse {
    return {
      unbondingDeposit: isSet(object.unbondingDeposit) ? UnbondingDeposit.fromJSON(object.unbondingDeposit) : undefined,
    };
  },

  toJSON(message: QueryGetUnbondingDepositResponse): unknown {
    const obj: any = {};
    message.unbondingDeposit !== undefined && (obj.unbondingDeposit = message.unbondingDeposit
      ? UnbondingDeposit.toJSON(message.unbondingDeposit)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetUnbondingDepositResponse>, I>>(
    object: I,
  ): QueryGetUnbondingDepositResponse {
    const message = createBaseQueryGetUnbondingDepositResponse();
    message.unbondingDeposit = (object.unbondingDeposit !== undefined && object.unbondingDeposit !== null)
      ? UnbondingDeposit.fromPartial(object.unbondingDeposit)
      : undefined;
    return message;
  },
};

function createBaseQueryAllUnbondingDepositRequest(): QueryAllUnbondingDepositRequest {
  return { pagination: undefined };
}

export const QueryAllUnbondingDepositRequest = {
  encode(message: QueryAllUnbondingDepositRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllUnbondingDepositRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllUnbondingDepositRequest();
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
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllUnbondingDepositRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllUnbondingDepositRequest>, I>>(
    object: I,
  ): QueryAllUnbondingDepositRequest {
    const message = createBaseQueryAllUnbondingDepositRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllUnbondingDepositResponse(): QueryAllUnbondingDepositResponse {
  return { unbondingDeposit: [], pagination: undefined };
}

export const QueryAllUnbondingDepositResponse = {
  encode(message: QueryAllUnbondingDepositResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.unbondingDeposit) {
      UnbondingDeposit.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllUnbondingDepositResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllUnbondingDepositResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.unbondingDeposit.push(UnbondingDeposit.decode(reader, reader.uint32()));
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
    return {
      unbondingDeposit: Array.isArray(object?.unbondingDeposit)
        ? object.unbondingDeposit.map((e: any) => UnbondingDeposit.fromJSON(e))
        : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllUnbondingDepositResponse): unknown {
    const obj: any = {};
    if (message.unbondingDeposit) {
      obj.unbondingDeposit = message.unbondingDeposit.map((e) => e ? UnbondingDeposit.toJSON(e) : undefined);
    } else {
      obj.unbondingDeposit = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllUnbondingDepositResponse>, I>>(
    object: I,
  ): QueryAllUnbondingDepositResponse {
    const message = createBaseQueryAllUnbondingDepositResponse();
    message.unbondingDeposit = object.unbondingDeposit?.map((e) => UnbondingDeposit.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetRefundPoolRequest(): QueryGetRefundPoolRequest {
  return { operatorAddress: "" };
}

export const QueryGetRefundPoolRequest = {
  encode(message: QueryGetRefundPoolRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.operatorAddress !== "") {
      writer.uint32(10).string(message.operatorAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetRefundPoolRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetRefundPoolRequest();
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

  fromJSON(object: any): QueryGetRefundPoolRequest {
    return { operatorAddress: isSet(object.operatorAddress) ? String(object.operatorAddress) : "" };
  },

  toJSON(message: QueryGetRefundPoolRequest): unknown {
    const obj: any = {};
    message.operatorAddress !== undefined && (obj.operatorAddress = message.operatorAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetRefundPoolRequest>, I>>(object: I): QueryGetRefundPoolRequest {
    const message = createBaseQueryGetRefundPoolRequest();
    message.operatorAddress = object.operatorAddress ?? "";
    return message;
  },
};

function createBaseQueryGetRefundPoolResponse(): QueryGetRefundPoolResponse {
  return { refundPool: undefined };
}

export const QueryGetRefundPoolResponse = {
  encode(message: QueryGetRefundPoolResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.refundPool !== undefined) {
      RefundPool.encode(message.refundPool, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetRefundPoolResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetRefundPoolResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.refundPool = RefundPool.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetRefundPoolResponse {
    return { refundPool: isSet(object.refundPool) ? RefundPool.fromJSON(object.refundPool) : undefined };
  },

  toJSON(message: QueryGetRefundPoolResponse): unknown {
    const obj: any = {};
    message.refundPool !== undefined
      && (obj.refundPool = message.refundPool ? RefundPool.toJSON(message.refundPool) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetRefundPoolResponse>, I>>(object: I): QueryGetRefundPoolResponse {
    const message = createBaseQueryGetRefundPoolResponse();
    message.refundPool = (object.refundPool !== undefined && object.refundPool !== null)
      ? RefundPool.fromPartial(object.refundPool)
      : undefined;
    return message;
  },
};

function createBaseQueryAllRefundPoolRequest(): QueryAllRefundPoolRequest {
  return { pagination: undefined };
}

export const QueryAllRefundPoolRequest = {
  encode(message: QueryAllRefundPoolRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRefundPoolRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRefundPoolRequest();
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

  fromJSON(object: any): QueryAllRefundPoolRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllRefundPoolRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRefundPoolRequest>, I>>(object: I): QueryAllRefundPoolRequest {
    const message = createBaseQueryAllRefundPoolRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllRefundPoolResponse(): QueryAllRefundPoolResponse {
  return { refundPool: [], pagination: undefined };
}

export const QueryAllRefundPoolResponse = {
  encode(message: QueryAllRefundPoolResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.refundPool) {
      RefundPool.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRefundPoolResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRefundPoolResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.refundPool.push(RefundPool.decode(reader, reader.uint32()));
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

  fromJSON(object: any): QueryAllRefundPoolResponse {
    return {
      refundPool: Array.isArray(object?.refundPool) ? object.refundPool.map((e: any) => RefundPool.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllRefundPoolResponse): unknown {
    const obj: any = {};
    if (message.refundPool) {
      obj.refundPool = message.refundPool.map((e) => e ? RefundPool.toJSON(e) : undefined);
    } else {
      obj.refundPool = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRefundPoolResponse>, I>>(object: I): QueryAllRefundPoolResponse {
    const message = createBaseQueryAllRefundPoolResponse();
    message.refundPool = object.refundPool?.map((e) => RefundPool.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetRefundRequest(): QueryGetRefundRequest {
  return { delegator: "", validator: "" };
}

export const QueryGetRefundRequest = {
  encode(message: QueryGetRefundRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.delegator !== "") {
      writer.uint32(10).string(message.delegator);
    }
    if (message.validator !== "") {
      writer.uint32(18).string(message.validator);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetRefundRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetRefundRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegator = reader.string();
          break;
        case 2:
          message.validator = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetRefundRequest {
    return {
      delegator: isSet(object.delegator) ? String(object.delegator) : "",
      validator: isSet(object.validator) ? String(object.validator) : "",
    };
  },

  toJSON(message: QueryGetRefundRequest): unknown {
    const obj: any = {};
    message.delegator !== undefined && (obj.delegator = message.delegator);
    message.validator !== undefined && (obj.validator = message.validator);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetRefundRequest>, I>>(object: I): QueryGetRefundRequest {
    const message = createBaseQueryGetRefundRequest();
    message.delegator = object.delegator ?? "";
    message.validator = object.validator ?? "";
    return message;
  },
};

function createBaseQueryGetRefundResponse(): QueryGetRefundResponse {
  return { refund: undefined };
}

export const QueryGetRefundResponse = {
  encode(message: QueryGetRefundResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.refund !== undefined) {
      Refund.encode(message.refund, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetRefundResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetRefundResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.refund = Refund.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetRefundResponse {
    return { refund: isSet(object.refund) ? Refund.fromJSON(object.refund) : undefined };
  },

  toJSON(message: QueryGetRefundResponse): unknown {
    const obj: any = {};
    message.refund !== undefined && (obj.refund = message.refund ? Refund.toJSON(message.refund) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetRefundResponse>, I>>(object: I): QueryGetRefundResponse {
    const message = createBaseQueryGetRefundResponse();
    message.refund = (object.refund !== undefined && object.refund !== null)
      ? Refund.fromPartial(object.refund)
      : undefined;
    return message;
  },
};

function createBaseQueryAllRefundRequest(): QueryAllRefundRequest {
  return { pagination: undefined };
}

export const QueryAllRefundRequest = {
  encode(message: QueryAllRefundRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRefundRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRefundRequest();
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

  fromJSON(object: any): QueryAllRefundRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllRefundRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRefundRequest>, I>>(object: I): QueryAllRefundRequest {
    const message = createBaseQueryAllRefundRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllRefundResponse(): QueryAllRefundResponse {
  return { refund: [], pagination: undefined };
}

export const QueryAllRefundResponse = {
  encode(message: QueryAllRefundResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.refund) {
      Refund.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRefundResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRefundResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.refund.push(Refund.decode(reader, reader.uint32()));
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

  fromJSON(object: any): QueryAllRefundResponse {
    return {
      refund: Array.isArray(object?.refund) ? object.refund.map((e: any) => Refund.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllRefundResponse): unknown {
    const obj: any = {};
    if (message.refund) {
      obj.refund = message.refund.map((e) => e ? Refund.toJSON(e) : undefined);
    } else {
      obj.refund = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRefundResponse>, I>>(object: I): QueryAllRefundResponse {
    const message = createBaseQueryAllRefundResponse();
    message.refund = object.refund?.map((e) => Refund.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
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
  DepositPool(request: QueryGetDepositPoolRequest): Promise<QueryGetDepositPoolResponse>;
  /** Queries a list of DepositPool items. */
  DepositPoolAll(request: QueryAllDepositPoolRequest): Promise<QueryAllDepositPoolResponse>;
  /** Queries a UnbondingDeposit by index. */
  UnbondingDeposit(request: QueryGetUnbondingDepositRequest): Promise<QueryGetUnbondingDepositResponse>;
  /** Queries a list of UnbondingDeposit items. */
  UnbondingDepositAll(request: QueryAllUnbondingDepositRequest): Promise<QueryAllUnbondingDepositResponse>;
  /** Queries a RefundPool by index. */
  RefundPool(request: QueryGetRefundPoolRequest): Promise<QueryGetRefundPoolResponse>;
  /** Queries a list of RefundPool items. */
  RefundPoolAll(request: QueryAllRefundPoolRequest): Promise<QueryAllRefundPoolResponse>;
  /** Queries a Refund by index. */
  Refund(request: QueryGetRefundRequest): Promise<QueryGetRefundResponse>;
  /** Queries a list of Refund items. */
  RefundAll(request: QueryAllRefundRequest): Promise<QueryAllRefundResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.Deposit = this.Deposit.bind(this);
    this.DepositAll = this.DepositAll.bind(this);
    this.DepositPool = this.DepositPool.bind(this);
    this.DepositPoolAll = this.DepositPoolAll.bind(this);
    this.UnbondingDeposit = this.UnbondingDeposit.bind(this);
    this.UnbondingDepositAll = this.UnbondingDepositAll.bind(this);
    this.RefundPool = this.RefundPool.bind(this);
    this.RefundPoolAll = this.RefundPoolAll.bind(this);
    this.Refund = this.Refund.bind(this);
    this.RefundAll = this.RefundAll.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  Deposit(request: QueryGetDepositRequest): Promise<QueryGetDepositResponse> {
    const data = QueryGetDepositRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "Deposit", data);
    return promise.then((data) => QueryGetDepositResponse.decode(new _m0.Reader(data)));
  }

  DepositAll(request: QueryAllDepositRequest): Promise<QueryAllDepositResponse> {
    const data = QueryAllDepositRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "DepositAll", data);
    return promise.then((data) => QueryAllDepositResponse.decode(new _m0.Reader(data)));
  }

  DepositPool(request: QueryGetDepositPoolRequest): Promise<QueryGetDepositPoolResponse> {
    const data = QueryGetDepositPoolRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "DepositPool", data);
    return promise.then((data) => QueryGetDepositPoolResponse.decode(new _m0.Reader(data)));
  }

  DepositPoolAll(request: QueryAllDepositPoolRequest): Promise<QueryAllDepositPoolResponse> {
    const data = QueryAllDepositPoolRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "DepositPoolAll", data);
    return promise.then((data) => QueryAllDepositPoolResponse.decode(new _m0.Reader(data)));
  }

  UnbondingDeposit(request: QueryGetUnbondingDepositRequest): Promise<QueryGetUnbondingDepositResponse> {
    const data = QueryGetUnbondingDepositRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "UnbondingDeposit", data);
    return promise.then((data) => QueryGetUnbondingDepositResponse.decode(new _m0.Reader(data)));
  }

  UnbondingDepositAll(request: QueryAllUnbondingDepositRequest): Promise<QueryAllUnbondingDepositResponse> {
    const data = QueryAllUnbondingDepositRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "UnbondingDepositAll", data);
    return promise.then((data) => QueryAllUnbondingDepositResponse.decode(new _m0.Reader(data)));
  }

  RefundPool(request: QueryGetRefundPoolRequest): Promise<QueryGetRefundPoolResponse> {
    const data = QueryGetRefundPoolRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "RefundPool", data);
    return promise.then((data) => QueryGetRefundPoolResponse.decode(new _m0.Reader(data)));
  }

  RefundPoolAll(request: QueryAllRefundPoolRequest): Promise<QueryAllRefundPoolResponse> {
    const data = QueryAllRefundPoolRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "RefundPoolAll", data);
    return promise.then((data) => QueryAllRefundPoolResponse.decode(new _m0.Reader(data)));
  }

  Refund(request: QueryGetRefundRequest): Promise<QueryGetRefundResponse> {
    const data = QueryGetRefundRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "Refund", data);
    return promise.then((data) => QueryGetRefundResponse.decode(new _m0.Reader(data)));
  }

  RefundAll(request: QueryAllRefundRequest): Promise<QueryAllRefundResponse> {
    const data = QueryAllRefundRequest.encode(request).finish();
    const promise = this.rpc.request("madeinblock.slashrefund.slashrefund.Query", "RefundAll", data);
    return promise.then((data) => QueryAllRefundResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
