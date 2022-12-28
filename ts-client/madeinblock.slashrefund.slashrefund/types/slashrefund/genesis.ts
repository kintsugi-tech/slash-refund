/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Deposit } from "./deposit";
import { DepositPool } from "./deposit_pool";
import { Params } from "./params";
import { Refund } from "./refund";
import { RefundPool } from "./refund_pool";
import { UnbondingDeposit } from "./unbonding_deposit";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

/** GenesisState defines the slashrefund module's genesis state. */
export interface GenesisState {
  params: Params | undefined;
  depositList: Deposit[];
  depositPoolList: DepositPool[];
  unbondingDepositList: UnbondingDeposit[];
  refundPoolList: RefundPool[];
  /** this line is used by starport scaffolding # genesis/proto/state */
  refundList: Refund[];
}

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    depositList: [],
    depositPoolList: [],
    unbondingDepositList: [],
    refundPoolList: [],
    refundList: [],
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.depositList) {
      Deposit.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.depositPoolList) {
      DepositPool.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.unbondingDepositList) {
      UnbondingDeposit.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.refundPoolList) {
      RefundPool.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.refundList) {
      Refund.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.depositList.push(Deposit.decode(reader, reader.uint32()));
          break;
        case 5:
          message.depositPoolList.push(DepositPool.decode(reader, reader.uint32()));
          break;
        case 6:
          message.unbondingDepositList.push(UnbondingDeposit.decode(reader, reader.uint32()));
          break;
        case 7:
          message.refundPoolList.push(RefundPool.decode(reader, reader.uint32()));
          break;
        case 8:
          message.refundList.push(Refund.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    return {
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined,
      depositList: Array.isArray(object?.depositList) ? object.depositList.map((e: any) => Deposit.fromJSON(e)) : [],
      depositPoolList: Array.isArray(object?.depositPoolList)
        ? object.depositPoolList.map((e: any) => DepositPool.fromJSON(e))
        : [],
      unbondingDepositList: Array.isArray(object?.unbondingDepositList)
        ? object.unbondingDepositList.map((e: any) => UnbondingDeposit.fromJSON(e))
        : [],
      refundPoolList: Array.isArray(object?.refundPoolList)
        ? object.refundPoolList.map((e: any) => RefundPool.fromJSON(e))
        : [],
      refundList: Array.isArray(object?.refundList) ? object.refundList.map((e: any) => Refund.fromJSON(e)) : [],
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    if (message.depositList) {
      obj.depositList = message.depositList.map((e) => e ? Deposit.toJSON(e) : undefined);
    } else {
      obj.depositList = [];
    }
    if (message.depositPoolList) {
      obj.depositPoolList = message.depositPoolList.map((e) => e ? DepositPool.toJSON(e) : undefined);
    } else {
      obj.depositPoolList = [];
    }
    if (message.unbondingDepositList) {
      obj.unbondingDepositList = message.unbondingDepositList.map((e) => e ? UnbondingDeposit.toJSON(e) : undefined);
    } else {
      obj.unbondingDepositList = [];
    }
    if (message.refundPoolList) {
      obj.refundPoolList = message.refundPoolList.map((e) => e ? RefundPool.toJSON(e) : undefined);
    } else {
      obj.refundPoolList = [];
    }
    if (message.refundList) {
      obj.refundList = message.refundList.map((e) => e ? Refund.toJSON(e) : undefined);
    } else {
      obj.refundList = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    message.depositList = object.depositList?.map((e) => Deposit.fromPartial(e)) || [];
    message.depositPoolList = object.depositPoolList?.map((e) => DepositPool.fromPartial(e)) || [];
    message.unbondingDepositList = object.unbondingDepositList?.map((e) => UnbondingDeposit.fromPartial(e)) || [];
    message.refundPoolList = object.refundPoolList?.map((e) => RefundPool.fromPartial(e)) || [];
    message.refundList = object.refundList?.map((e) => Refund.fromPartial(e)) || [];
    return message;
  },
};

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
