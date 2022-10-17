/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

export interface Deposit {
  depositorAddress: string;
  validatorAddress: string;
  shares: string;
}

const baseDeposit: object = {
  depositorAddress: "",
  validatorAddress: "",
  shares: "",
};

export const Deposit = {
  encode(message: Deposit, writer: Writer = Writer.create()): Writer {
    if (message.depositorAddress !== "") {
      writer.uint32(10).string(message.depositorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    if (message.shares !== "") {
      writer.uint32(26).string(message.shares);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Deposit {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDeposit } as Deposit;
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
          message.shares = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Deposit {
    const message = { ...baseDeposit } as Deposit;
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
    if (object.shares !== undefined && object.shares !== null) {
      message.shares = String(object.shares);
    } else {
      message.shares = "";
    }
    return message;
  },

  toJSON(message: Deposit): unknown {
    const obj: any = {};
    message.depositorAddress !== undefined &&
      (obj.depositorAddress = message.depositorAddress);
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    message.shares !== undefined && (obj.shares = message.shares);
    return obj;
  },

  fromPartial(object: DeepPartial<Deposit>): Deposit {
    const message = { ...baseDeposit } as Deposit;
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
    if (object.shares !== undefined && object.shares !== null) {
      message.shares = object.shares;
    } else {
      message.shares = "";
    }
    return message;
  },
};

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
