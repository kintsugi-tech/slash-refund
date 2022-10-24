/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

export interface UnbondingDeposit {
  delegatorAddress: string;
  validatorAddress: string;
  unbondingDepositEntry: string;
}

const baseUnbondingDeposit: object = {
  delegatorAddress: "",
  validatorAddress: "",
  unbondingDepositEntry: "",
};

export const UnbondingDeposit = {
  encode(message: UnbondingDeposit, writer: Writer = Writer.create()): Writer {
    if (message.delegatorAddress !== "") {
      writer.uint32(10).string(message.delegatorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    if (message.unbondingDepositEntry !== "") {
      writer.uint32(26).string(message.unbondingDepositEntry);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): UnbondingDeposit {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUnbondingDeposit } as UnbondingDeposit;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegatorAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        case 3:
          message.unbondingDepositEntry = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnbondingDeposit {
    const message = { ...baseUnbondingDeposit } as UnbondingDeposit;
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = String(object.delegatorAddress);
    } else {
      message.delegatorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = String(object.validatorAddress);
    } else {
      message.validatorAddress = "";
    }
    if (
      object.unbondingDepositEntry !== undefined &&
      object.unbondingDepositEntry !== null
    ) {
      message.unbondingDepositEntry = String(object.unbondingDepositEntry);
    } else {
      message.unbondingDepositEntry = "";
    }
    return message;
  },

  toJSON(message: UnbondingDeposit): unknown {
    const obj: any = {};
    message.delegatorAddress !== undefined &&
      (obj.delegatorAddress = message.delegatorAddress);
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    message.unbondingDepositEntry !== undefined &&
      (obj.unbondingDepositEntry = message.unbondingDepositEntry);
    return obj;
  },

  fromPartial(object: DeepPartial<UnbondingDeposit>): UnbondingDeposit {
    const message = { ...baseUnbondingDeposit } as UnbondingDeposit;
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = object.delegatorAddress;
    } else {
      message.delegatorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = object.validatorAddress;
    } else {
      message.validatorAddress = "";
    }
    if (
      object.unbondingDepositEntry !== undefined &&
      object.unbondingDepositEntry !== null
    ) {
      message.unbondingDepositEntry = object.unbondingDepositEntry;
    } else {
      message.unbondingDepositEntry = "";
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
