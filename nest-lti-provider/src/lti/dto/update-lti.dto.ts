import { PartialType } from '@nestjs/mapped-types';
import { CreateLtiDto } from './create-lti.dto';

export class UpdateLtiDto extends PartialType(CreateLtiDto) {}
