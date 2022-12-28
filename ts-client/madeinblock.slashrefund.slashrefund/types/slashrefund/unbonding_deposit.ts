/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { UnbondingDepositEntry } from "./unbonding_deposit_entry";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

/**
 * option (gogoproto.goproto_getters)  = false;
 * option (gogoproto.equal)            = false;
 * option (gogoproto.goproto_stringer) = false;
 */
export interface UnbondingDeposit {
  depositorAddress: string;
  validatorAddress: string;
  entries: UnbondingDepositEntry[];
}

function createBaseUnbondingDeposit(): UnbondingDeposit {
  return { depositorAddress: "", validatorAddress: "", entries: [] };
}

export const UnbondingDeposit = {
  encode(message: UnbondingDeposit, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.depositorAddress !== "") {
      writer.uint32(10).string(message.depositorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    for (const v of message.entries) {
      UnbondingDepositEntry.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnbondingDeposit {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnbondingDeposit();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.depositorAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        case 3:
          message.entries.push(UnbondingDepositEntry.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnbondingDeposit {
    return {
      depositorAddress: isSet(object.depositorAddress) ? String(object.depositorAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
      entries: Array.isArray(object?.entries) ? object.entries.map((e: any) => UnbondingDepositEntry.fromJSON(e)) : [],
    };
  },

  toJSON(message: UnbondingDeposit): unknown {
    const obj: any = {};
    message.depositorAddress !== undefined && (obj.depositorAddress = message.depositorAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    if (message.entries) {
      obj.entries = message.entries.map((e) => e ? UnbondingDepositEntry.toJSON(e) : undefined);
    } else {
      obj.entries = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UnbondingDeposit>, I>>(object: I): UnbondingDeposit {
    const message = createBaseUnbondingDeposit();
    message.depositorAddress = object.depositorAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    message.entries = object.entries?.map((e) => UnbondingDepositEntry.fromPartial(e)) || [];
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
