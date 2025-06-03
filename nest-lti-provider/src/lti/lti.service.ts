import { Injectable } from '@nestjs/common';
import { CreateLtiDto } from './dto/create-lti.dto';
import { UpdateLtiDto } from './dto/update-lti.dto';

@Injectable()
export class LtiService {
  create(createLtiDto: CreateLtiDto) {
    return 'This action adds a new lti';
  }

  findAll() {
    return `This action returns all lti`;
  }

  findOne(id: number) {
    return `This action returns a #${id} lti`;
  }

  update(id: number, updateLtiDto: UpdateLtiDto) {
    return `This action updates a #${id} lti`;
  }

  remove(id: number) {
    return `This action removes a #${id} lti`;
  }
}
