# Copyright © 2022 Durudex
#
# This file is part of Durudex: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# Durudex is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Durudex. If not, see <https://www.gnu.org/licenses/>.

FROM golang:1.18 AS builder

RUN go version

COPY . /github.com/durudex/discord-promo-bot/
WORKDIR /github.com/durudex/discord-promo-bot/

RUN go mod download
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bot ./cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/durudex/discord-promo-bot/bot .

CMD ["./bot"]
